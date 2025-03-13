INSERT INTO
  users (name, email)
VALUES
  ('Alice Smith', 'alice@example.com'),
  ('Bob Jones', 'bob@example.com'),
  ('Carol Wilson', 'carol@example.com');

INSERT INTO
  products (name, description, price)
VALUES
  (
    'Premium Coffee Subscription',
    'Artisanal coffee delivered monthly',
    19.99
  ),
  (
    'Standard Coffee Subscription',
    'Great quality coffee delivered monthly',
    14.99
  ),
  (
    'Premium Tea Subscription',
    'Exotic tea selection delivered monthly',
    16.99
  ),
  (
    'Tea Sampler Subscription',
    'Try different teas each month',
    12.99
  );

INSERT INTO
  subscriptions (
    user_id,
    product_id,
    start_date,
    next_billing_date,
    status
  )
VALUES
  (
    1,
    1,
    CURRENT_DATE - INTERVAL '15 days',
    CURRENT_DATE + INTERVAL '15 days',
    'active'
  ),
  (
    2,
    2,
    CURRENT_DATE - INTERVAL '40 days',
    CURRENT_DATE - INTERVAL '10 days',
    'hold'
  ),
  (
    3,
    3,
    CURRENT_DATE - INTERVAL '31 days',
    CURRENT_DATE - INTERVAL '1 days',
    'active'
  );

INSERT INTO
  bills (
    subscription_id,
    amount,
    status,
    created_at,
    paid_at
  )
VALUES
  (
    1,
    14.99,
    'paid',
    CURRENT_DATE - INTERVAL '15 days',
    CURRENT_DATE - INTERVAL '15 days'
  );

INSERT INTO
  bills (subscription_id, amount, status, created_at)
VALUES
  (
    2,
    19.99,
    'pending',
    CURRENT_DATE - INTERVAL '10 day'
  );
