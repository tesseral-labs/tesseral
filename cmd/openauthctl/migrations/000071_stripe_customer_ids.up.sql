alter table projects add column stripe_customer_id varchar;
create index on projects (stripe_customer_id);
