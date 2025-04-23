CREATE TABLE reception (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    execution_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    pvz_id UUID NOT NULL REFERENCES pvz(id) ON DELETE CASCADE,
    status VARCHAR(50) NOT NULL CHECK (status IN ('in_progress', 'close'))
);