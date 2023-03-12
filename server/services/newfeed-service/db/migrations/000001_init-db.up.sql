CREATE TABLE IF NOT EXISTS newfeed.users(
    id uuid PRIMARY KEY,
    email text
);
CREATE TABLE IF NOT EXISTS newfeed.postfeeds_by_user_id(
    user_id uuid,
    author_id uuid,
    post_id uuid,
    created_at timestamp,
    PRIMARY KEY (user_id, created_at)
)
WITH CLUSTERING ORDER BY (created_at DESC);
CREATE TABLE IF NOT EXISTS newfeed.storyfeeds_by_user_id(
    user_id uuid,
    author_id uuid,
    story_id uuid,
    created_at timestamp,
    PRIMARY KEY (user_id, created_at)
)
WITH CLUSTERING ORDER BY (created_at DESC);
CREATE TABLE IF NOT EXISTS newfeed.petifyfeeds_by_user_id(
    user_id uuid,
    author_id uuid,
    petify_id uuid,
    created_at timestamp,
    PRIMARY KEY (user_id, created_at)
)
WITH CLUSTERING ORDER BY (created_at DESC);
CREATE TABLE IF NOT EXISTS newfeed.postfeeds_by_user_id_and_post_id(
    user_id uuid,
    author_id uuid,
    post_id uuid,
    created_at timestamp,
    PRIMARY KEY ((user_id, post_id))
);
CREATE TABLE IF NOT EXISTS newfeed.storyfeeds_by_user_id_and_story_id(
    user_id uuid,
    author_id uuid,
    story_id uuid,
    created_at timestamp,
    PRIMARY KEY ((user_id, story_id))
);
CREATE TABLE IF NOT EXISTS newfeed.petifyfeeds_by_user_id_and_petify_id(
    user_id uuid,
    author_id uuid,
    petify_id uuid,
    created_at timestamp,
    PRIMARY KEY ((user_id, petify_id))
);
