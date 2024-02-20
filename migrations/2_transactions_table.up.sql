CREATE TABLE IF NOT EXISTS transacoes (
    id SERIAL PRIMARY KEY,
    id_cliente INT,
    valor DECIMAL(20, 2) NOT NULL,
    tipo CHAR(1) NOT NULL CHECK (tipo IN ('c', 'd')),
    descricao TEXT NOT NULL,
    realizada_em TEXT NOT NULL
);