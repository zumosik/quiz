package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"rest_grpc/pb/auth"
	"rest_grpc/pb/files"
	"rest_grpc/utils"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	maxMBSize = 100 // 100 mb
)

type contextType int

const (
	ctxTokenKey contextType = iota
)

type Server struct {
	r          *chi.Mux
	l          *slog.Logger
	fCl        files.FilesClient
	aCl        auth.AuthClient
	CtxTimeout time.Duration
}

func (s *Server) Run(addr string) error {
	s.l.Info(fmt.Sprintf("Starting server on addr: %s", addr))
	return http.ListenAndServe(addr, s.r)
}

func MustNew(l *slog.Logger, filesAddr string, authAddr string, timeout time.Duration) *Server {
	r := chi.NewRouter()

	var srv Server

	srv.configureRouter(r)

	conn, err := grpc.Dial(filesAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		l.Error(utils.WrapErr("error in grpc.Dial", err))
		os.Exit(1)
	}

	conn2, err := grpc.Dial(authAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		l.Error(utils.WrapErr("error in grpc.Dial", err))
		os.Exit(1)
	}

	srv.fCl = files.NewFilesClient(conn)
	srv.aCl = auth.NewAuthClient(conn2)

	srv.l = l
	srv.r = r

	srv.CtxTimeout = timeout

	return &srv
}

func (s *Server) configureRouter(r *chi.Mux) {
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/users", func(r chi.Router) {
		r.Use(middleware.AllowContentType("application/json"))

		r.Post("/login", s.Login())
		r.Post("/register", s.Register())
	})

	r.Route("/files", func(r chi.Router) {
		r.Use(middleware.AllowContentType("multipart/form-data", "application/json"))

		r.Post("/upload", s.Upload())
	})
}

func (s *Server) Upload() http.HandlerFunc {
	const requestMultipartFormFileName = "file"

	type request struct {
		Name      string `json:"name" validate:"required"`
		CreatedAt string `json:"created_at" validate:"required"` // RFC3339
	}

	type response struct {
		Response
		ID string `json:"id"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), s.CtxTimeout)
		defer cancel()

		err := r.ParseMultipartForm(maxMBSize << 20) // 10 MB max file size
		if err != nil {
			render.JSON(w, r, response{
				Response: Response{
					StatusCode: 400,
					Ok:         "",
					Error:      "cant get multipart form",
				},
			})

			w.WriteHeader(400)
		}

		file, _, err := r.FormFile(requestMultipartFormFileName)
		if err != nil {

			render.JSON(w, r, response{
				Response: Response{
					StatusCode: 400,
					Ok:         "",
					Error:      "cant get your file",
				},
			})

			w.WriteHeader(400)
		}
		defer file.Close()

		bytes, err := io.ReadAll(file)
		if err != nil {

			render.JSON(w, r, response{
				Response: Response{
					StatusCode: 500,
					Ok:         "",
					Error:      "cant read file",
				},
			})

			w.WriteHeader(500) // idk if this 500 or 400
		}

		var req request
		json.NewDecoder(r.Body).Decode(&req)
		err = validator.New().Struct(req)
		if err != nil {
			render.JSON(w, r, response{
				Response: Response{
					StatusCode: 400,
					Ok:         "",
					Error:      "invalid data",
				},
			})

			w.WriteHeader(400)
		}

		// getting id from middleware
		userID := r.Context().Value(ctxTokenKey).(string)

		t, err := time.Parse(req.CreatedAt, time.RFC3339)
		if err != nil {
			render.JSON(w, r, response{
				Response: Response{
					StatusCode: 400,
					Ok:         "",
					Error:      "you must pass created_at field in RFC3339 format",
				},
			})

			w.WriteHeader(400)
		}

		res, err := s.fCl.UploadFile(ctx, &files.UploadFileRequest{
			UserId: userID,
			File: &files.File{
				Content:   bytes,
				Name:      req.Name,
				CreatedAt: timestamppb.New(t),
			},
		})

		if err != nil {
			render.JSON(w, r, response{
				Response: Response{
					StatusCode: 500,
					Ok:         "",
					Error:      "internal error",
				},
			})

			w.WriteHeader(500)
		}

		render.JSON(w, r, response{
			Response: Response{
				StatusCode: 200,
				Ok:         "ok",
			},
			ID: res.GetFile().GetId(),
		})

		w.WriteHeader(200)
	}
}

func (s *Server) Login() http.HandlerFunc {
	type request struct {
		Email    string `json:"email" validate:"required"`
		Passowrd string `json:"password" validate:"required,email"`
		AppID    int    `json:"app_id"`
	}

	type response struct {
		Response
		Token string `json:"token"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), s.CtxTimeout)
		defer cancel()

		var req request
		json.NewDecoder(r.Body).Decode(&req)
		err := validator.New().Struct(req)
		if err != nil {
			render.JSON(w, r, response{
				Response: Response{
					StatusCode: 400,
					Ok:         "",
					Error:      "invalid data",
				},
				Token: "",
			})

			w.WriteHeader(400)
		}
		res, err := s.aCl.Login(ctx, &auth.LoginRequest{
			Email:    req.Email,
			Password: req.Passowrd,
			AppId:    int32(req.AppID),
		})

		if err != nil {
			render.JSON(w, r, response{
				Response: Response{
					StatusCode: 500,
					Ok:         "",
					Error:      "internal server error",
				},
				Token: "",
			})

			w.WriteHeader(500)
		}

		render.JSON(w, r, response{
			Response: Response{
				StatusCode: 200,
				Ok:         "ok",
			},
			Token: res.GetToken(),
		})
	}
}

func (s *Server) Register() http.HandlerFunc {
	type request struct {
		Email    string `json:"email" validate:"required"`
		Passowrd string `json:"password" validate:"required,email"`
		AppID    int    `json:"app_id"`
	}

	type response struct {
		Response
		Token string `json:"token"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), s.CtxTimeout)
		defer cancel()

		var req request
		json.NewDecoder(r.Body).Decode(&req)
		err := validator.New().Struct(req)
		if err != nil {
			render.JSON(w, r, response{
				Response: Response{
					StatusCode: 400,
					Ok:         "",
					Error:      "invalid data",
				},
				Token: "",
			})

			w.WriteHeader(400)
		}

		// we doesnt need any info from resp so just checking error
		_, err = s.aCl.Register(ctx, &auth.RegisterRequest{
			Email:    req.Email,
			Password: req.Passowrd,
		})

		if err != nil {
			render.JSON(w, r, response{
				Response: Response{
					StatusCode: 500,
					Ok:         "",
					Error:      "internal server error",
				},
				Token: "",
			})

			w.WriteHeader(500)
		}

		res2, err := s.aCl.Login(ctx, &auth.LoginRequest{
			Email:    req.Email,
			Password: req.Passowrd,
			AppId:    int32(req.AppID),
		})

		if err != nil {
			render.JSON(w, r, response{
				Response: Response{
					StatusCode: 500,
					Ok:         "",
					Error:      "internal server error",
				},
				Token: "",
			})

			w.WriteHeader(500)
		}

		render.JSON(w, r, response{
			Response: Response{
				StatusCode: 200,
				Ok:         "ok",
			},
			Token: res2.GetToken(),
		})
	}
}

func (s *Server) GetIDFromToken(next http.Handler) http.Handler {
	type response struct {
		Response
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), s.CtxTimeout)
		defer cancel()

		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			render.JSON(w, r, response{
				Response: Response{
					StatusCode: http.StatusUnauthorized,
					Ok:         "",
					Error:      "you need to pass token",
				},
			})

			w.WriteHeader(http.StatusUnauthorized)
		}
		tokenString = tokenString[len("Bearer "):]

		res, err := s.aCl.GetID(ctx, &auth.GetIDRequest{Token: tokenString})
		if err != nil {
			render.JSON(w, r, response{
				Response: Response{
					StatusCode: http.StatusUnauthorized,
					Ok:         "",
					Error:      "invalid token",
				},
			})

			w.WriteHeader(http.StatusUnauthorized)
		}

		newCtx := context.WithValue(r.Context(), ctxTokenKey, res.GetUserId())

		next.ServeHTTP(w, r.WithContext(newCtx))
	})
}
