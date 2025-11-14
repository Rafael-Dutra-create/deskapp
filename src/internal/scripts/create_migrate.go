package main

import (
	"fmt"
	"os"
	"time"
)


// CreateMigration cria uma nova migração com os arquivos up e down
func CreateMigration(migrationsPath, name string) error {
    // Usar timestamp no formato YYYYMMDDhhmmss em vez de Unix epoch
    timestamp := time.Now().Format("20060102150405")
    
    // Nome do arquivo com timestamp no formato correto
    upFilename := fmt.Sprintf("%s/%s_%s.up.sql", migrationsPath, timestamp, name)
    downFilename := fmt.Sprintf("%s/%s_%s.down.sql", migrationsPath, timestamp, name)
    
    // Criar arquivo UP
    upFile, err := os.Create(upFilename)
    if err != nil {
        return err
    }
    defer upFile.Close()
    
    // Adicionar comentário com metadados no arquivo SQL
    upFile.WriteString(fmt.Sprintf("-- Migration: %s\n-- Created at: %s\n-- Created by: %s\n\n", 
        name, time.Now().Format(time.RFC3339), os.Getenv("USER")))
    
    // Criar arquivo DOWN
    downFile, err := os.Create(downFilename)
    if err != nil {
        return err
    }
    defer downFile.Close()
    
    downFile.WriteString(fmt.Sprintf("-- Migration: %s\n-- Created at: %s\n-- Created by: %s\n\n", 
        name, time.Now().Format(time.RFC3339), os.Getenv("USER")))
    
    fmt.Printf("✅ Migração criada: %s\n", name)
    fmt.Printf("📁 Arquivos: %s, %s\n", upFilename, downFilename)
    
    return nil
}