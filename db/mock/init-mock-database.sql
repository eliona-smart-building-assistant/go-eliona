--  This file is part of the eliona project.
--  Copyright © 2022 LEICOM iTEC AG. All Rights Reserved.
--  ______ _ _
-- |  ____| (_)
-- | |__  | |_  ___  _ __   __ _
-- |  __| | | |/ _ \| '_ \ / _` |
-- | |____| | | (_) | | | | (_| |
-- |______|_|_|\___/|_| |_|\__,_|
--
--  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING
--  BUT NOT LIMITED  TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
--  NON INFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
--  DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
--  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

-- Use this file to initialize a mocking database. The database have to be PostgreSQL.
-- You can use any cloud service or a docker container to create a local database. An example
-- docker-compose.yml file is also provided in this directory.

create schema if not exists public;

create table if not exists public.asset_type
(
    asset_type         text not null primary key,
    custom             boolean default true not null,
    payload_fct        text,
    vendor             text,
    model              text,
    translation        jsonb,
    urldoc             text,
    allowed_inactivity interval,
    iv_asset_type      integer,
    icon               text
);

create table if not exists public.asset
(
    asset_id    serial primary key,
    proj_id     text,
    gai         text not null,
    name        text,
    device_pkey text unique,
    asset_type  text,
    lat         double precision,
    lon         double precision,
    storey      smallint,
    description text,
    tags        text[],
    ar          boolean default false not null,
    tracker     boolean default false not null,
    loc_ref     integer,
    func_ref    integer,
    urldoc      text,
    unique (gai, proj_id)
    );

create table if not exists public.attribute_schema
(
    id              serial primary key,
    asset_type      text                     not null,
    attribute_type  text,
    attribute       text                     not null,
    subtype         text    default ''::text not null,
    enable          boolean default true     not null,
    translation     jsonb,
    unit            text,
    formula         text,
    scale           numeric,
    zero            double precision,
    precision       smallint,
    min             numeric,
    max             numeric,
    step            numeric,
    map             json,
    pipeline_mode   text,
    pipeline_raster text[],
    viewer          boolean default false    not null,
    ar              boolean default false    not null,
    seq             smallint,
    source_path     text[],
    virtual         boolean,
    unique (asset_type, subtype, attribute)
    );

create table if not exists public.heap
(
    asset_id            integer                                not null,
    subtype             text                                   not null,
    his                 boolean                  default true  not null,
    ts                  timestamp with time zone default now() not null,
    data                jsonb,
    valid               boolean,
    allowed_inactivity  interval,
    update_cnt          bigint                   default 1     not null,
    update_cnt_reset_ts timestamp with time zone default now() not null,
    primary key (asset_id, subtype)
    );

create table if not exists public.eliona_app (
    app_name    text primary key,
    category    text,
    active      boolean default false,
    initialised boolean default false
);

create schema if not exists versioning;

create table if not exists versioning.patches
(
    app_name    text                                   not null,
    patch_name  text                                   not null,
    applied_tsz timestamp with time zone default now() not null,
    applied_by  text                                   not null,
    requires    text[],
    conflicts   text[],
    primary key (app_name, patch_name)
    );

