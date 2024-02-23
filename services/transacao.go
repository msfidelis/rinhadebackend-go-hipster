package services

import (
	"context"
	"database/sql"
	"fmt"
	"main/entities"
	"main/pkg/database"

	"github.com/uptrace/bun"
)

const (
	OperacaoCredito = "c"
	OperacaoDebito  = "d"
)

func BuscaCliente(ctx context.Context, db *bun.DB, id string) (*entities.Cliente, error) {
	functionName := "BuscaCLiente"
	cliente := new(entities.Cliente)

	err := db.NewSelect().Model(cliente).Where("id_cliente = ?", id).Scan(ctx)
	if err != nil {
		fmt.Printf("[%s] Erro ao encontrar o cliente %v:\n", functionName, err)
		return cliente, err
	}
	return cliente, nil
}

func BuscaClienteTx(ctx context.Context, tx bun.Tx, id string) (*entities.Cliente, error) {
	functionName := "BuscaCLiente"
	cliente := new(entities.Cliente)

	err := tx.NewSelect().Model(cliente).Where("id_cliente = ?", id).Scan(ctx)
	if err != nil {
		fmt.Printf("[%s] Erro ao encontrar o cliente %v:\n", functionName, err)
		tx.Rollback()
		return cliente, err
	}
	return cliente, nil
}

func Crebito(transacao entities.Transacao) (novoSaldo int64, limite int64, inconsistente bool, err error) {
	functionName := fmt.Sprintf("OperacaoDe%s", transacao.Tipo)

	ctx := context.Background()
	db := database.GetDB()

	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		fmt.Printf("[%s] Erro ao iniciar a transação: %v\n", functionName, err)
		return
	}

	// Garantir que um rollback será feito em caso de erro
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	cliente, err := BuscaClienteTx(ctx, tx, transacao.IDCliente)
	if err != nil {
		fmt.Printf("[%s] Erro ao encontrar o cliente %v:\n", functionName, err)
		return
	}

	// Cálculo do novo saldo com base no tipo de operação
	switch transacao.Tipo {
	case OperacaoCredito:
		novoSaldo = cliente.Saldo + transacao.Valor
	case OperacaoDebito:
		novoSaldo = cliente.Saldo - transacao.Valor
		if novoSaldo < -cliente.Limite {
			inconsistente = true
			err = fmt.Errorf("[%s] Operação excederia o limite do cliente", functionName)
			return 0, 0, true, err
		}
	default:
		err = fmt.Errorf("[%s] Tipo de operação inválido", functionName)
		return 0, 0, false, err
	}

	// Atualizar saldo do cliente
	// _, err = tx.ExecContext(ctx, "UPDATE clientes SET saldo = ? WHERE id_cliente = ?", novoSaldo, transacao.IDCliente)
	// if err != nil {
	// 	fmt.Printf("[%s] Erro ao atualizar o saldo do cliente: %v\n", functionName, err)
	// 	return
	// }

	_, err = tx.NewUpdate().Model((*entities.Cliente)(nil)).
		Set("saldo = ?", novoSaldo).
		Where("id_cliente = ?", transacao.IDCliente).
		Exec(ctx)

	if err != nil {
		fmt.Printf("[%s] Erro ao atualizar o saldo do cliente: %v\n", functionName, err)
		return 0, 0, false, err
	}

	_, err = tx.NewInsert().
		Model(&transacao).
		Exec(ctx)

	// Inserir transação
	// _, err = tx.ExecContext(ctx, "INSERT INTO transacoes (id_cliente, valor, tipo, descricao) VALUES (?, ?, ?, ?)", transacao.IDCliente, transacao.Valor, tipo, transacao.Descricao)
	if err != nil {
		fmt.Printf("[%s] Erro ao inserir a transação: %v\n", functionName, err)
		return 0, 0, false, err
	}

	// Commit da Transação
	err = tx.Commit()
	if err != nil {
		fmt.Printf("[%s] Erro ao fazer commit da transação: %v\n", functionName, err)
		return 0, 0, false, err
	}

	limite = cliente.Limite
	return novoSaldo, cliente.Limite, inconsistente, nil
}
