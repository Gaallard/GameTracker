package service

import (
	"errors"
	"gametracker/db"
	"gametracker/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

var jwtSecret = []byte("gametracker_secret_key_2024") // En producción usar variable de entorno

// Register crea un nuevo usuario
func (s *AuthService) Register(req models.RegisterRequest) (*models.User, error) {
	// Verificar si el usuario ya existe (solo verificar existencia, no cargar datos)
	var count int64
	if err := db.DB.Model(&models.User{}).Where("username = ? OR email = ?", req.Username, req.Email).Count(&count).Error; err != nil {
		return nil, errors.New("error al verificar usuario existente")
	}
	if count > 0 {
		return nil, errors.New("usuario o email ya existe")
	}

	// Crear nuevo usuario
	user := &models.User{
		Username:  req.Username,
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Encriptar contraseña
	if err := user.HashPassword(req.Password); err != nil {
		return nil, errors.New("error al encriptar contraseña")
	}

	// Guardar en base de datos
	if err := db.DB.Create(user).Error; err != nil {
		return nil, errors.New("error al crear usuario")
	}

	// Limpiar contraseña antes de devolver
	user.Password = ""
	return user, nil
}

// Login autentica un usuario
func (s *AuthService) Login(req models.LoginRequest) (*models.AuthResponse, error) {
	var user models.User
	
	// Buscar usuario por username o email
	if err := db.DB.Where("username = ? OR email = ?", req.Username, req.Username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("usuario no encontrado")
		}
		return nil, errors.New("error al buscar usuario")
	}

	// Verificar contraseña
	if !user.CheckPassword(req.Password) {
		return nil, errors.New("contraseña incorrecta")
	}

	// Generar token JWT
	token, err := s.generateToken(user.ID, user.Username)
	if err != nil {
		return nil, errors.New("error al generar token")
	}

	// Limpiar contraseña antes de devolver
	user.Password = ""

	return &models.AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

// generateToken genera un token JWT
func (s *AuthService) generateToken(userID uint, username string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 días
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ValidateToken valida un token JWT
func (s *AuthService) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("método de firma inválido")
		}
		return jwtSecret, nil
	})
}

// GetUserFromToken extrae información del usuario desde el token
func (s *AuthService) GetUserFromToken(token *jwt.Token) (uint, string, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, "", errors.New("token inválido")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, "", errors.New("user_id no encontrado en token")
	}

	username, ok := claims["username"].(string)
	if !ok {
		return 0, "", errors.New("username no encontrado en token")
	}

	return uint(userID), username, nil
}
