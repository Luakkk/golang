create table if not exists users (
  id serial primary key,
  name varchar(255) not null,
  email varchar(255) not null unique,
  age int not null check (age >= 0),
  created_at timestamptz not null default now()
);

insert into users (name, email, age) values ('John Doe', 'john@example.com', 30)
on conflict (email) do nothing;