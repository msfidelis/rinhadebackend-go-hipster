package services

import (
	"context"
	"database/sql"
	"fmt"
	"main/entities"
	"main/pkg/database"
)

func Credito(transacao entities.Transacao) (float64, float64, bool, error) {
	functionName := "OperacaoDeCredito"

	ctx := context.Background()
	db := database.GetDB()

	errChan := make(chan error, 2)
	doneChan := make(chan bool, 2)

	tx, err := db.BeginTx(ctx, &sql.TxOptions{})

	cliente, err := BuscaCliente(transacao.IDCliente)
	if err != nil {
		fmt.Printf("[%s] Erro ao encontrar o cliente %v:\n", functionName, err)
		return 0, 0, false, err
	}
	novoSaldo := cliente.Saldo - transacao.Valor

	go func() {
		_, err = tx.NewUpdate().Model((*entities.Cliente)(nil)).
			Set("saldo = + ?", novoSaldo).
			Where("id_cliente = ?", transacao.IDCliente).
			Exec(ctx)
		errChan <- err
		doneChan <- true
	}()

	go func() {
		_, err = tx.NewInsert().
			Model(&transacao).
			Exec(ctx)

		errChan <- err
		doneChan <- true
	}()

	// Aguardar a conclusão das duas operações
	<-doneChan
	<-doneChan

	// Verificar erros
	close(errChan)
	for e := range errChan {
		if e != nil {
			tx.Rollback()
			return 0, 0, false, e
		}
	}

	// Commit da Transação
	if err = tx.Commit(); err != nil {
		fmt.Printf("[%s] Erro ao fazer commit da transação: %v\n", functionName, err)
		return 0, 0, false, err
	}
	return novoSaldo, cliente.Limite, false, nil
}

func Debito(transacao entities.Transacao) (float64, float64, bool, error) {
	functionName := "OperacaoDeDebito"

	ctx := context.Background()
	db := database.GetDB()

	errChan := make(chan error, 2)
	doneChan := make(chan bool, 2)

	tx, err := db.BeginTx(ctx, &sql.TxOptions{})

	cliente, err := BuscaCliente(transacao.IDCliente)
	if err != nil {
		fmt.Printf("[%s] Erro ao encontrar o cliente %v:\n", functionName, err)
		return 0, 0, false, err
	}

	novoSaldo := cliente.Saldo - transacao.Valor
	if novoSaldo < -cliente.Limite {
		tx.Rollback()
		return cliente.Saldo, cliente.Limite, true, fmt.Errorf("[%s] Operação excederia o limite do cliente", functionName)
	}

	go func() {
		_, err = tx.NewUpdate().Model((*entities.Cliente)(nil)).
			Set("saldo = ?", novoSaldo).
			Where("id_cliente = ?", transacao.IDCliente).
			Exec(ctx)
		errChan <- err
		doneChan <- true
	}()

	go func() {
		_, err = tx.NewInsert().
			Model(&transacao).
			Exec(ctx)

		errChan <- err
		doneChan <- true
	}()

	// Aguardar a conclusão das duas operações
	<-doneChan
	<-doneChan

	// Verificar erros
	close(errChan)
	for e := range errChan {
		if e != nil {
			tx.Rollback()
			return 0, 0, false, e
		}
	}

	// Commit da Transação
	if err = tx.Commit(); err != nil {
		fmt.Printf("[%s] Erro ao fazer commit da transação: %v\n", functionName, err)
		return 0, 0, false, err
	}

	return novoSaldo, cliente.Limite, false, nil
}

func BuscaCliente(id string) (*entities.Cliente, error) {
	functionName := "BuscaCLiente"
	cliente := new(entities.Cliente)

	ctx := context.Background()
	db := database.GetDB()

	err := db.NewSelect().Model(cliente).Where("id_cliente = ?", id).Scan(ctx)
	if err != nil {
		fmt.Printf("[%s] Erro ao encontrar o cliente %v:\n", functionName, err)
		return cliente, err
	}
	return cliente, nil
}
