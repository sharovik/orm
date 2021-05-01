create table table_name
(
	res_id integer not null
		constraint table_name_history_id_fk
			references history,
	id integer
		constraint table_name_pk
			primary key
);

create unique index table_name_id_uindex
	on table_name (id);

