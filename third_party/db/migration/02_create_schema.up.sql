CREATE SCHEMA bibi;

SET search_path to bibi ;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE bibi.skin_types (
    id uuid DEFAULT uuid_generate_v4(),
    name character varying(255) NOT NULL,
    created_time timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_time timestamp with time zone DEFAULT NULL,

    CONSTRAINT skin_types_pkey PRIMARY KEY (id)
);

CREATE TABLE bibi.skin_problems (
    id uuid DEFAULT uuid_generate_v4(),
    name character varying(255) NOT NULL,
    created_time timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_time timestamp with time zone DEFAULT NULL,

    CONSTRAINT skin_problems_pkey PRIMARY KEY (id)
);

CREATE TABLE bibi.users (
    id uuid DEFAULT uuid_generate_v4(),
    email character varying(255) NOT NULL,
    password_hash character varying(255) NOT NULL,
    gender character varying(255), --- male, female, lgbtq+
    birthdate date DEFAULT NULL,
    skin_type_id uuid DEFAULT NULL,
    skin_problem_id uuid DEFAULT NULL,
    created_time timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT users_pkey PRIMARY KEY (id),
    CONSTRAINT users_email_unique UNIQUE (email),
    CONSTRAINT users_skin_type_id_fkey FOREIGN KEY (skin_type_id)
        REFERENCES skin_types (id) MATCH FULL
        ON UPDATE CASCADE,
    CONSTRAINT users_skin_problem_id_fkey FOREIGN KEY (skin_problem_id)
        REFERENCES skin_problems (id) MATCH FULL
        ON UPDATE CASCADE
);

CREATE TABLE bibi.admins (
    id uuid DEFAULT uuid_generate_v4(),
    username character varying(255) NOT NULL,
    password_hash character varying(255) NOT NULL,
    created_time timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT admins_pkey PRIMARY KEY (id)
);

INSERT INTO bibi.admins(username,password_hash) VALUES ('admin','e1f5a62b304f8cec9e3505b6ec793f18f874aa898d81dda9e9b2916f7a47ae05.7d24c8d7fca782be87c84d1421578a56f51254edf67bec063ff3e95185410862');

CREATE TABLE bibi.product_types (
    id uuid DEFAULT uuid_generate_v4(),
    name character varying(255) NOT NULL,
    created_time timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_time timestamp with time zone DEFAULT NULL,

    CONSTRAINT product_types_pkey PRIMARY KEY (id)
);

CREATE TABLE bibi.product_categories (
    id uuid DEFAULT uuid_generate_v4(),
    name character varying(255) NOT NULL,
    created_time timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_time timestamp with time zone DEFAULT NULL,

    CONSTRAINT product_categories_pkey PRIMARY KEY (id)
);

CREATE TABLE bibi.countries (
    id uuid DEFAULT uuid_generate_v4(),
    name character varying(255) NOT NULL,
    created_time timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_time timestamp with time zone DEFAULT NULL,

    CONSTRAINT countries_pkey PRIMARY KEY (id)
);

CREATE TABLE bibi.banners (
    id uuid DEFAULT uuid_generate_v4(),
    name character varying(255) NOT NULL,
    area_code character varying(255) NOT NULL,
    created_time timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT banners_pkey PRIMARY KEY (id),
    CONSTRAINT banners_area_code_unique UNIQUE (area_code)
);


CREATE TABLE bibi.banner_images (
    id uuid DEFAULT uuid_generate_v4(),
    banner_id uuid DEFAULT NULL,
    path character varying(255) NOT NULL,
    created_time timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT banner_images_pkey PRIMARY KEY (id),
    CONSTRAINT banner_images_banner_id_fkey FOREIGN KEY (banner_id)
        REFERENCES banners (id) MATCH FULL
        ON UPDATE CASCADE
        ON DELETE CASCADE
);

CREATE TABLE bibi.products (
    id uuid DEFAULT uuid_generate_v4(),
    brand character varying(255),
    name character varying(255),
    short_description text,
    description text,
    size character varying(255),
    price decimal(10,2),
    product_type_id uuid DEFAULT NULL,
    product_category_id uuid DEFAULT NULL,
    skin_type_id uuid DEFAULT NULL,
    country_id uuid DEFAULT NULL,
    tags text[] DEFAULT '{}',
    created_time timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_time timestamp with time zone DEFAULT NULL,

    CONSTRAINT products_pkey PRIMARY KEY (id),
    CONSTRAINT products_product_type_id_fkey FOREIGN KEY (product_type_id)
        REFERENCES product_types (id) MATCH FULL
        ON UPDATE CASCADE,
    CONSTRAINT products_product_category_id_fkey FOREIGN KEY (product_category_id)
        REFERENCES product_categories (id) MATCH FULL
        ON UPDATE CASCADE,
    CONSTRAINT products_skin_type_id_fkey FOREIGN KEY (skin_type_id)
        REFERENCES skin_types (id) MATCH FULL
        ON UPDATE CASCADE,
    CONSTRAINT products_country_id_fkey FOREIGN KEY (country_id)
        REFERENCES countries (id) MATCH FULL
        ON UPDATE CASCADE
);

CREATE TABLE bibi.product_images (
    id uuid DEFAULT uuid_generate_v4(),
    product_id uuid DEFAULT NULL,
    path character varying(255) NOT NULL,
    created_time timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT product_images_pkey PRIMARY KEY (id),
    CONSTRAINT product_images_product_id_fkey FOREIGN KEY (product_id)
        REFERENCES products (id) MATCH FULL
        ON UPDATE CASCADE
        ON DELETE CASCADE
);

CREATE TABLE bibi.product_recommends (
    id uuid DEFAULT uuid_generate_v4(),
    product_id uuid DEFAULT NULL,
    created_time timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT product_recommends_pkey PRIMARY KEY (id),
    CONSTRAINT product_recommends_product_id_fkey FOREIGN KEY (product_id)
        REFERENCES products (id) MATCH FULL
        ON UPDATE CASCADE
        ON DELETE CASCADE
);