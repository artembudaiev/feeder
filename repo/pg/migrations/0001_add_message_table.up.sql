create table if not exists message (
    id UUID default gen_random_uuid() primary key,
    text varchar(50)
    );