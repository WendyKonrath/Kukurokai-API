package routes

import (
	config "go-api/db"
	"go-api/middleware"
	"go-api/models"
	"go-api/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func SetupProductRoutes(app *fiber.App) {
	produtoGroup := app.Group("/produtos", middleware.JWTMiddleware())

	produtoGroup.Get("/", GetProdutos)
	produtoGroup.Get("/:id", GetProduto)
	produtoGroup.Post("/", CreateProduto)
	produtoGroup.Put("/:id", EditProduto)
	produtoGroup.Delete("/:id", DelProdutos)
}

func GetProdutos(c *fiber.Ctx) error {
	user := c.Locals("user").(jwt.MapClaims)
	role := user["role"].(string)

	if role != "admin" && role != "superadmin" {
		return c.Status(403).JSON(fiber.Map{"error": "Acesso proibido"})
	}

	var produtos []models.Produto
	if err := config.DB.Find(&produtos).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Erro ao buscar produtos"})
	}

	return c.JSON(produtos)
}

func GetProduto(c *fiber.Ctx) error {
	user := c.Locals("user").(jwt.MapClaims)
	role := user["role"].(string)

	if role != "admin" && role != "superadmin" {
		return c.Status(403).JSON(fiber.Map{"error": "Acesso proibido"})
	}

	id := c.Params("id")
	var produto models.Produto

	if err := config.DB.First(&produto, "id = ?", id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Produto não encontrado"})
	}

	// Verificar o tipo do produto e carregar informações específicas
	switch produto.Tipo {
	case models.Fisico:
		var produtoFisico models.ProdutoFisico
		if err := config.DB.Where("id = ?", id).First(&produtoFisico).Error; err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Detalhes do produto físico não encontrados"})
		}
		return c.JSON(produtoFisico)

	case models.Servico:
		var produtoServico models.ProdutoServico
		if err := config.DB.Where("id = ?", id).First(&produtoServico).Error; err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Detalhes do serviço não encontrados"})
		}
		return c.JSON(produtoServico)

	default:
		return c.JSON(produto)
	}
}

func CreateProduto(c *fiber.Ctx) error {
	user := c.Locals("user").(jwt.MapClaims)
	role := user["role"].(string)

	if role != "admin" && role != "superadmin" {
		return c.Status(403).JSON(fiber.Map{"error": "Acesso proibido"})
	}

	type ProdutoRequest struct {
		Produto models.Produto         `json:"produto"`
		Fisico  *models.ProdutoFisico  `json:"fisico,omitempty"`
		Servico *models.ProdutoServico `json:"servico,omitempty"`
	}

	var req ProdutoRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Requisição inválida"})
	}

	// Validação dos dados do produto
	if err := utils.Validate.Struct(req.Produto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Dados inválidos", "details": err.Error()})
	}

	// Criar o produto base
	if err := config.DB.Create(&req.Produto).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Erro ao criar produto"})
	}

	// Criar detalhes específicos baseado no tipo
	switch req.Produto.Tipo {
	case models.Fisico:
		if req.Fisico == nil {
			return c.Status(400).JSON(fiber.Map{"error": "Dados do produto físico são obrigatórios"})
		}
		req.Fisico.Produto = req.Produto
		if err := utils.Validate.Struct(req.Fisico); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Dados inválidos do produto físico", "details": err.Error()})
		}
		if err := config.DB.Create(req.Fisico).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Erro ao criar produto físico"})
		}
		return c.JSON(req.Fisico)

	case models.Servico:
		if req.Servico == nil {
			return c.Status(400).JSON(fiber.Map{"error": "Dados do serviço são obrigatórios"})
		}
		req.Servico.Produto = req.Produto
		if err := utils.Validate.Struct(req.Servico); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Dados inválidos do serviço", "details": err.Error()})
		}
		if err := config.DB.Create(req.Servico).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Erro ao criar serviço"})
		}
		return c.JSON(req.Servico)

	default:
		return c.JSON(req.Produto)
	}
}

func EditProduto(c *fiber.Ctx) error {
	user := c.Locals("user").(jwt.MapClaims)
	role := user["role"].(string)

	if role != "admin" && role != "superadmin" {
		return c.Status(403).JSON(fiber.Map{"error": "Acesso proibido"})
	}

	id := c.Params("id")
	var produto models.Produto
	if err := config.DB.First(&produto, "id = ?", id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Produto não encontrado"})
	}

	type ProdutoRequest struct {
		Produto models.Produto         `json:"produto"`
		Fisico  *models.ProdutoFisico  `json:"fisico,omitempty"`
		Servico *models.ProdutoServico `json:"servico,omitempty"`
	}

	var req ProdutoRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Requisição inválida"})
	}

	// Atualizar produto base
	if err := utils.Validate.Struct(req.Produto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Dados inválidos", "details": err.Error()})
	}

	product := req.Produto
	product.ID = id // Garantir que o ID não seja alterado
	if err := config.DB.Save(&product).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Erro ao atualizar produto"})
	}

	// Atualizar detalhes específicos baseado no tipo
	switch product.Tipo {
	case models.Fisico:
		if req.Fisico == nil {
			return c.Status(400).JSON(fiber.Map{"error": "Dados do produto físico são obrigatórios"})
		}
		req.Fisico.Produto = product
		if err := utils.Validate.Struct(req.Fisico); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Dados inválidos do produto físico", "details": err.Error()})
		}
		if err := config.DB.Where("id = ?", id).Updates(req.Fisico).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Erro ao atualizar produto físico"})
		}
		return c.JSON(req.Fisico)

	case models.Servico:
		if req.Servico == nil {
			return c.Status(400).JSON(fiber.Map{"error": "Dados do serviço são obrigatórios"})
		}
		req.Servico.Produto = product
		if err := utils.Validate.Struct(req.Servico); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Dados inválidos do serviço", "details": err.Error()})
		}
		if err := config.DB.Where("id = ?", id).Updates(req.Servico).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Erro ao atualizar serviço"})
		}
		return c.JSON(req.Servico)

	default:
		return c.JSON(product)
	}
}

func DelProdutos(c *fiber.Ctx) error {
	user := c.Locals("user").(jwt.MapClaims)
	role := user["role"].(string)

	if role != "admin" && role != "superadmin" {
		return c.Status(403).JSON(fiber.Map{"error": "Acesso proibido"})
	}

	id := c.Params("id")
	var produto models.Produto
	if err := config.DB.First(&produto, "id = ?", id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Produto não encontrado"})
	}

	// Deletar o produto e seus detalhes específicos
	if err := config.DB.Delete(&produto).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Erro ao deletar produto"})
	}

	return c.SendStatus(204)
}
