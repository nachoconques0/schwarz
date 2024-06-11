BEGIN;

CREATE TABLE schwarz.coupon (
  id UUID PRIMARY KEY,
  name TEXT NOT NULL,
  amount FLOAT NOT NULL,
  used BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE schwarz.shopping_cart (
  id UUID PRIMARY KEY,
  coupon_id UUID,
  items JSONB DEFAULT NULL,
  amount FLOAT NOT NULL,
  total FLOAT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL
);

-- Update triggers
CREATE TRIGGER set_updated_at
  BEFORE INSERT OR UPDATE ON schwarz.coupon
  FOR EACH ROW
  EXECUTE PROCEDURE schwarz.set_updated_at ();

CREATE TRIGGER set_updated_at
  BEFORE INSERT OR UPDATE ON schwarz.shopping_cart
  FOR EACH ROW
  EXECUTE PROCEDURE schwarz.set_updated_at ();

COMMIT;
