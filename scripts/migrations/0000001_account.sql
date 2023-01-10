create table if not exists public.account
(
    id      integer generated always as identity
        constraint account_pk
            primary key,
    balance numeric default 0 not null,
    version integer default 0 not null
);