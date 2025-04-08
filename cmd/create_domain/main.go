package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func toPascalCase(name string) string {
	return cases.Title(language.English).String(name)
}

/*
func toPlural(name string) string {
	if strings.HasSuffix(name, "s") {
		return name + "es"
	}
	return name + "s"
}
*/

func parseFields(args []string) (string, string) {
	if len(args) == 0 {
		return "    ID string `db:\"id\" json:\"id\"`\n    Field string `db:\"field\" json:\"field\"`", "\"id\", \"field\""
	}

	var structFields []string
	var fields []string

	structFields = append(structFields, "    ID string `db:\"id\" json:\"id\"`")
	fields = append(fields, "\"id\"")
	for _, field := range args {
		parts := strings.Split(field, ":")
		if len(parts) != 2 {
			log.Fatalf("Campo inválido: %s. Use Nome:Tipo", field)
		}
		name := parts[0]
		typ := parts[1]
		dbTag := strings.ToLower(name)
		fieldLine := fmt.Sprintf("    %s %s `db:\"%s\" json:\"%s\"`", name, typ, dbTag, dbTag)
		structFields = append(structFields, fieldLine)
		fields = append(fields, "\""+dbTag+"\"")
	}

	structFields = append(structFields, "    CreatedAt string `db:\"created_at\" json:\"created_at\"`")
	structFields = append(structFields, "    UpdatedAt string `db:\"updated_at\" json:\"updated_at\"`")

	fields = append(fields, "\"created_at\"")
	fields = append(fields, "\"updated_at\"")

	return strings.Join(structFields, "\n"), strings.Join(fields, ", ")
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Você precisa passar o nome do domínio. Ex: go run main.go car [Nome:string Idade:int]")
	}

	domain := strings.ToLower(os.Args[1])
	structName := toPascalCase(domain)
	//plural := toPlural(domain)

	structFields, fields := parseFields(os.Args[2:])

	modelPath := fmt.Sprintf("model/%s.go", domain)
	modelContent := fmt.Sprintf(
		"package model\n\ntype %s struct {\n%s\n}\n\nvar %sFields = []string{%s}",
		structName,
		structFields,
		structName,
		fields,
	)

	err := os.WriteFile(modelPath, []byte(modelContent), 0644)
	if err != nil {
		log.Fatalf("Erro ao criar model: %v", err)
	}
	fmt.Printf("✅ Model gerado: %s\n", modelPath)

	registryPath := "util/registry.go"
	registryContent, err := os.ReadFile(registryPath)
	if err != nil {
		log.Fatalf("Erro ao ler registry.go: %v", err)
	}

	importLine := "\n\t\"api_boilerplate/model\""
	if !strings.Contains(string(registryContent), importLine) {
		registryContent = []byte(strings.Replace(string(registryContent), "import (", "import "+importLine, 1))
	}

	registerLine := fmt.Sprintf("\tRegisterGenericResource[model.%s](r, db, \"%s\", model.%sFields)", structName, domain, structName)
	if strings.Contains(string(registryContent), registerLine) {
		fmt.Println("ℹ️  Já registrado em registry.go")
	} else {
		lines := strings.Split(string(registryContent), "\n")
		var updated []string
		inserted := false
		for _, line := range lines {
			updated = append(updated, line)
			if strings.Contains(line, "func RegisterDomains") && !inserted {
				updated = append(updated, registerLine)
				inserted = true
			}
		}

		err = os.WriteFile(registryPath, []byte(strings.Join(updated, "\n")), 0644)
		if err != nil {
			log.Fatalf("Erro ao atualizar registry.go: %v", err)
		}

		fmt.Println("✅ Registro adicionado em registry.go")
	}
}
