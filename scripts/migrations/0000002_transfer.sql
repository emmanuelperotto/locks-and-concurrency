create table if not exists public.transfer
(
    id              integer generated always as identity
        constraint transfer_pk
            primary key,
    amount          numeric not null
        constraint non_negative
            check (amount > (0)::numeric),
    from_account_id integer not null
        constraint from_account___fk
            references public.account,
    to_account_id   integer not null
        constraint to_account___fk
            references public.account
);

