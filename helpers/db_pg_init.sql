drop table if exists orders;
create table orders (id text, content text, is_final bool);

drop table if exists accounting;
create table accounting_orders (order_id text, price int, paid_at timestamp);

drop table if exists prices;
create table prices (order_id text, price int);

drop table if exists shipping;
create table shipping (order_id text, shipped_at timestamp);