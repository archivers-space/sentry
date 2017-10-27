package core

// create collections table if it doesn't exist
const qCollectionCreateTable = `
CREATE TABLE IF NOT EXISTS collections (
  id               UUID PRIMARY KEY,
  created          timestamp NOT NULL,
  updated          timestamp NOT NULL,
  creator          text NOT NULL DEFAULT '',
  title            text NOT NULL DEFAULT '',
  url              text NOT NULL DEFAULT ''
);`

// list collections by reverse cronological date created
// paginated
const qCollections = `
SELECT
  id, created, updated, creator, title, description, url
FROM collections 
ORDER BY created DESC 
LIMIT $1 OFFSET $2;`

// list collections by creator
const qCollectionsByCreator = `
SELECT 
  id, created, updated, creator, title, description, url 
FROM collections
WHERE creator = $4
ORDER BY $3
LIMIT $1 OFFSET $2;`

// check for existence of a collection
const qCollectionExists = `
SELECT exists(SELECT 1 FROM collections WHERE id = $1)
`

// insert a collection
const qCollectionInsert = `
INSERT INTO collections 
  (id, created, updated, creator, title, description, url ) 
VALUES ($1, $2, $3, $4, $5, $6, $7);`

// update an existing collection, selecting by ID
const qCollectionUpdate = `
UPDATE collections 
SET created=$2, updated=$3, creator=$4, title=$5, description=$6, url=$7
WHERE id = $1;`

// read collection info by ID
const qCollectionById = `
SELECT 
  id, created, updated, creator, title, description, url 
FROM collections 
WHERE id = $1;`

const qCollectionByUrl = `
SELECT
  id, created, updated, creator, title, description, url 
FROM collections 
WHERE url = $1;`

// deleted a collection
const qCollectionDelete = `
DELETE from collections 
WHERE id = $1;`

const qCollectionItemCreateTable = `
CREATE TABLE IF NOT EXISTS collection_items (
  collection_id    UUID NOT NULL,
  url_id           text NOT NULL default '',
  index            integer NOT NULL default -1,
  description      text NOT NULL default '',
  PRIMARY KEY      (collection_id, url_id)
);`

const qCollectionItemInsert = `
INSERT INTO collection_items
  (collection_id, url_id, index, description)
VALUES
  ($1, $2, $3, $4);`

const qCollectionItemUpdate = `
UPDATE collection_items
SET index = $3, description = $4
WHERE collection_id = $1 and url_id = $2;`

const qCollectionItemDelete = `
DELETE FROM collection_items 
WHERE collection_id = $1 AND url_id = $2;`

const qCollectionItemExists = `
SELECT exists(SELECT 1 FROM collection_items where collection_id = $1 AND url_id = $2);`

const qCollectionItemById = `
SELECT
  ci.collection_id, u.id, u.hash, u.url, u.title, ci.index, ci.description
FROM collection_items as ci, urls as u
WHERE collection_id = $1 AND url_id = $2 AND u.id = ci.url_id;`

const qCollectionLength = `
SELECT count(1) FROM collection_items WHERE collection_id = $1;`

const qCollectionItems = `
SELECT
  ci.collection_id, u.id, u.hash, u.url, u.title, ci.index, ci.description
FROM collection_items as ci, urls as u
WHERE collection_id = $4
AND u.id = ci.url_id
ORDER BY $3
LIMIT $1 OFFSET $2;`

// insert a dataRepo
const qDataRepoInsert = `
INSERT INTO data_repos 
  (id, created, updated, title, description, url) 
VALUES ($1, $2, $3, $4, $5, $6);`

// update an existing dataRepo, selecting by ID
const qDataRepoUpdate = `
UPDATE data_repos 
SET created=$2, updated=$3, title=$4, description=$5, url=$6
WHERE id = $1;`

const qDataRepoCreateTable = `
CREATE TABLE IF NOT EXISTS data_repos (
  id               UUID PRIMARY KEY NOT NULL,
  created          timestamp NOT NULL default (now() at time zone 'utc'),
  updated          timestamp NOT NULL default (now() at time zone 'utc'),
  title            text NOT NULL default '',
  description      text NOT NULL default '',
  url              text NOT NULL default '',
  deleted          boolean default false
);`

const qDataRepoExists = `SELECT exists(SELECT 1 FROM data_repos WHERE id = $1)`

// read dataRepo info by ID
const qDataRepoById = `
SELECT 
  id, created, updated, title, description, url 
FROM data_repos 
WHERE id = $1;`

// deleted a dataRepo
const qDataRepoDelete = `
DELETE from data_repos 
WHERE id = $1;`

// list data_repos by reverse cronological date created
// paginated
const qDataRepos = `
SELECT
  id, created, updated, title, description, url
FROM data_repos 
ORDER BY created DESC 
LIMIT $1 OFFSET $2;`

// create links table
const qLinkCreateTable = `
CREATE TABLE IF NOT EXISTS links (
  created          timestamp NOT NULL,
  updated          timestamp NOT NULL,
  src              text NOT NULL references urls(url) ON DELETE CASCADE,
  dst              text NOT NULL references urls(url) ON DELETE CASCADE,
  PRIMARY KEY      (src, dst)
);`

// check to see if a link exists
const qLinkExists = `SELECT exists(select 1 from links where src = $1 and dst = $2)`

// read a link
const qLinkRead = `
SELECT
  created, updated, src, dst
FROM links
WHERE
  src = $1 AND
  dst = $2;`

// insert a link
const qLinkInsert = `
INSERT INTO LINKS
  (created, updated, src, dst)
VALUES
  ($1, $2, $3, $4);`

// update a link
const qLinkUpdate = `
UPDATE links SET
  created = $1, updated = $2
WHERE
  src = $3 AND
  dst = $4`

// delete a link
const qLinkDelete = `
DELETE FROM links
WHERE
  src = $1 AND
  dst = $2;`

const qMetadataCreateTable = `
CREATE TABLE IF NOT EXISTS metadata (
  hash             text NOT NULL default '',
  time_stamp       timestamp NOT NULL,
  key_id           text NOT NULL default '',
  subject          text NOT NULL,
  prev             text NOT NULL default '',
  meta             json,
  deleted          boolean default false
);
`

// list latest metadata entries by reverse cronological order
// paginated
const qMetadataLatest = `
SELECT
  hash, time_stamp, key_id, subject, prev, meta 
FROM metadata 
WHERE 
  key_id = $1 AND 
  subject = $2 
ORDER BY time_stamp DESC;`

// select metadata entries for a given subject hash
// TODO - should this be paginated?
const qMetadataForSubject = `
SELECT
  hash, time_stamp, key_id, subject, prev, meta  
FROM metadata
WHERE 
  subject = $1 AND 
  deleted = false AND 
  meta IS NOT NULL;`

// count the number of metatadata entries for a a given user key
// omitting empty content
const qMetadataCountForKey = `
SELECT
  count(1)
FROM metadata
WHERE
  -- confirm is not empty hash
  hash != '1220e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855' AND
  key_id = $1;`

// pagniated selection of only the latest metadata entry for each user-key/subject pair
// paginated
const qMetadataLatestForKey = `
SELECT DISTINCT ON (subject)
  hash, time_stamp, key_id, subject, prev, meta
FROM metadata
WHERE
  key_id = $1 and
  deleted = false
ORDER BY subject, time_stamp DESC
LIMIT $2 OFFSET $3;`

// insert a metdata entry
const qMetadataInsert = `
INSERT INTO metadata
  (hash, time_stamp, key_id, subject, prev, meta, deleted)
VALUES 
  ($1, $2, $3, $4, $5, $6, false);`

const qPrimerCreateTable = `
CREATE TABLE IF NOT EXISTS primers (
  id               UUID PRIMARY KEY NOT NULL,
  created          timestamp NOT NULL default (now() at time zone 'utc'),
  updated          timestamp NOT NULL default (now() at time zone 'utc'),
  short_title      text NOT NULL default '',
  title            text NOT NULL default '',
  description      text NOT NULL default '',
  parent_id        text NOT NULL default '', -- this should be "UUID references primers(id)", but then we'd need to accept null values, no bueno
  stats            json,
  meta             json,
  deleted          boolean default false
);`

const qPrimerExists = `SELECT exists(SELECT 1 FROM primers WHERE id = $1)`

// read a primer for a given Id
const qPrimerById = `
SELECT
  id, created, updated, short_title, title, description, 
  parent_id, stats, meta
FROM primers 
WHERE 
  deleted = false AND
  id = $1;`

// insert a primer
const qPrimerInsert = `
INSERT INTO primers
  (id, created, updated, short_title, title, description, parent_id, stats, meta)
VALUES
  ($1, $2, $3, $4, $5, $6, $7, $8, $9);`

// update an existing primer
const qPrimerUpdate = `
UPDATE PRIMERS set
  created = $2, updated = $3, short_title = $4, title = $5, description = $6,
  parent_id = $7, stats = $8, meta = $9
WHERE
  deleted = false AND
  id = $1;`

// 'delete' a primer
const qPrimerDelete = `
UPDATE primers SET
  deleted = true
WHERE id = $1;`

// read children for a given primer ID. children only, decendants are omitted.
const qPrimerSubPrimers = `
SELECT
  id, created, updated, short_title, title, description, 
  parent_id, stats, meta
FROM primers
WHERE 
  deleted = false AND
  parent_id = $1;`

// read all of a primer's sources
const qPrimerSources = `
SELECT
  id, created, updated, title, description, url, primer_id, crawl, stale_duration,
  last_alert_sent, meta, stats
FROM sources
WHERE
  deleted = false AND
  primer_id = $1;`

// enumerate primers
const qPrimersCount = `SELECT count(1) FROM primers WHERE deleted = false`

// list primers by reverse chronolgical created date, no hierarchy is observed
// paginated
const qPrimersList = `
SELECT
  id, created, updated, short_title, title, description,
  parent_id, stats, meta
FROM primers
WHERE
  deleted = false
ORDER BY created DESC
LIMIT $1 OFFSET $2;`

// list primers that have no parent ("root" or "base" primers), order by reverse chronological created date.
// paginated
const qBasePrimersList = `
select
  id, created, updated, short_title, title, description,
  parent_id, stats, meta
from primers
where
  deleted = false and
  parent_id = ''
order by created desc
limit $1 offset $2;`

const qSourceCreateTable = `
CREATE TABLE IF NOT EXISTS sources (
  id               UUID PRIMARY KEY NOT NULL,
  created          timestamp NOT NULL default (now() at time zone 'utc'),
  updated          timestamp NOT NULL default (now() at time zone 'utc'),
  title            text NOT NULL default '',
  description      text NOT NULL default '',
  url              text UNIQUE NOT NULL,
  primer_id        UUID references primers(id) ON DELETE CASCADE,
  crawl            boolean default true,
  stale_duration   integer NOT NULL DEFAULT 43200000, -- defaults to 12 hours, column needs to be multiplied by 1000000 to become a poper duration
  last_alert_sent  timestamp,
  stats            json,
  meta             json,
  deleted          boolean default false
);`

const qSourceExists = `SELECT exists(SELECT 1 FROM sources WHERE id = $1)`
const qSourceExistsByUrl = `SELECT exists(SELECT 1 FROM sources WHERE url = $1)`

// select
const qSourcesCount = `SELECT count(1) FROM sources;`

// list sources, ordered by reverse chronological created date
// paginated
const qSourcesList = `
SELECT
  id, created, updated, title, description, url, primer_id, crawl, stale_duration,
  last_alert_sent, meta, stats
FROM sources
WHERE 
  deleted = false
ORDER BY created DESC
LIMIT $1 OFFSET $2;`

// read a source for a given ID
const qSourceById = `
SELECT
  id, created, updated, title, description, url, primer_id, crawl, stale_duration,
  last_alert_sent, meta, stats
FROM sources 
WHERE
  deleted = false AND
  id = $1;`

// read a source for a given url string
const qSourceByUrl = `
SELECT
  id, created, updated, title, description, url, primer_id, crawl, stale_duration,
  last_alert_sent, meta, stats
FROM sources 
WHERE
  deleted = false AND
  url = $1;`

// add a source
const qSourceInsert = `
INSERT INTO sources 
  (id, created, updated, title, description, url, primer_id, crawl, stale_duration,
   last_alert_sent, meta, stats) 
VALUES 
  ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);`

// update an exsiting source
const qSourceUpdate = `
UPDATE sources 
SET 
  created = $2, updated = $3, title = $4, description = $5, url = $6, primer_id = $7, 
  crawl = $8, stale_duration = $9, last_alert_sent = $10, meta = $11, stats = $12
WHERE
  deleted = false AND
  id = $1;`

// 'delete' a source
const qSourceDelete = `
UPDATE sources
SET
  deleted = true
WHERE 
  url = $1;`

// count how many URLs this source matches from the urls table
// the passed in argument can take the form '%[arg]%' to ignore position within
// the url string
// warning - big & slow
const qSourceUrlCount = `
SELECT count(1) 
FROM urls 
WHERE 
  url ilike $1;`

// list sources that have crawl column set to true, ordered by reverse chronolgical created date
// paginated
const qSourcesCrawling = `
SELECT
  id, created, updated, title, description, url, primer_id, crawl, stale_duration,
  last_alert_sent, meta, stats
FROM sources
WHERE
  deleted = false and
  crawl = true
ORDER BY created DESC
LIMIT $1 OFFSET $2;`

// enumerate all urls for a given source that are suspected to have content
// this generates the "content urls" stat.
const qSourceContentUrlCount = `
SELECT count(1) 
FROM urls 
WHERE
  hash != '' AND
  hash != '1220e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855' AND
  content_sniff != 'text/html; charset=utf-8' AND
  url ilike $1;`

// enumerate all urls for a given source that are suspected to have content *and*
// have an entry in the metadata table
// this generates the "content urls" stat.
const qSourceContentWithMetadataCount = `
SELECT count(1)
FROM urls 
WHERE 
  url ilike $1 AND 
  content_sniff != 'text/html; charset=utf-8' AND
  -- confirm is not empty hash
  hash != '1220e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855' AND
  hash != '' AND
  EXISTS (SELECT null FROM metadata WHERE urls.hash = metadata.subject);`

const qSourceUndescribedContentUrls = `
select
  url, created, updated, last_head, last_get, status, content_type, content_sniff,
  content_length, file_name, title, id, headers_took, download_took, headers, meta, hash
from urls 
where 
  url ilike $1
  and content_sniff != 'text/html; charset=utf-8'
  and last_get is not null
  -- confirm is not empty hash
  and hash != '1220e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855'
  and not exists (select null from metadata where urls.hash = metadata.subject) 
limit $2 offset $3;`

const qSourceDescribedContentUrls = `
select
  url, created, updated, last_head, last_get, status, content_type, content_sniff,
  content_length, file_name, title, id, headers_took, download_took, headers, meta, hash
from urls 
where 
  url ilike $1
  and content_sniff != 'text/html; charset=utf-8'
  and last_get is not null
  -- confirm is not empty hash
  and hash != '1220e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855'
  and exists (select null from metadata where urls.hash = metadata.subject) 
limit $2 offset $3;`

const qSnapshotCreateTable = `
CREATE TABLE IF NOT EXISTS snapshots (
  url              text NOT NULL references urls(url) ON DELETE CASCADE,
  created          timestamp NOT NULL,
  status           integer NOT NULL DEFAULT 0,
  duration         integer NOT NULL DEFAULT 0,
  meta             json,
  hash             text NOT NULL DEFAULT ''
);`

const qSnapeshotExists = `SELECT exists(SELECT 1 FROM snapshots WHERE hash = $1)`

const qSnapshotsByUrl = `
select
  url, created, status, duration, meta, hash
from snapshots 
where 
  url = $1;`

const qSnapshotInsert = `
insert into snapshots 
  (url, created, status, duration, meta, hash)
values 
  ($1, $2, $3, $4, $5, $6);`

const qUrlsSearch = `
select
  url, created, updated, last_head, last_get, status, content_type, content_sniff,
  content_length, file_name, title, id, headers_took, download_took, headers, meta, hash
from urls 
where 
  url ilike $1 
limit $2 offset $3;`

const qUrlsCreateTable = `
CREATE TABLE IF NOT EXISTS urls (
  url              text PRIMARY KEY NOT NULL,
  created          timestamp NOT NULL,
  updated          timestamp NOT NULL,
  last_head        timestamp,
  last_get         timestamp,
  status           integer NOT NULL default 0,
  content_type     text NOT NULL default '',
  content_sniff    text NOT NULL default '',
  content_length   bigint NOT NULL default 0,
  file_name        text NOT NULL default '',
  title            text NOT NULL default '',
  id               text NOT NULL default '',
  headers_took     integer NOT NULL default 0,
  download_took    integer NOT NULL default 0,
  headers          json,
  meta             json,
  hash             text NOT NULL default ''
);`

const qUrlExists = `SELECT exists(SELECT 1 FROM urls WHERE id = $1)`

const qUrlsList = `
select
  url, created, updated, last_head, last_get, status, content_type, content_sniff,
  content_length, file_name, title, id, headers_took, download_took, headers, meta, hash
from urls 
order by created desc 
limit $1 offset $2;`

const qContentUrlsList = `
select
  url, created, updated, last_head, last_get, status, content_type, content_sniff,
  content_length, file_name, title, id, headers_took, download_took, headers, meta, hash
from urls 
where
  last_get is not null and
  content_sniff != 'text/html; charset=utf-8' and
  content_sniff != '' and
  hash != ''
order by created desc
limit $1 offset $2;`

const qContentUrlsCount = `
select
  count(1)
from urls 
where
  last_get is not null and
  content_sniff != 'text/html; charset=utf-8' and
  content_sniff != '' and
  hash != ''
`

const qUrlsFetched = `
select
  url, created, updated, last_head, last_get, status, content_type, content_sniff,
  content_length, file_name, title, id, headers_took, download_took, headers, meta, hash
from urls 
where
  last_get is not null 
order by created desc
limit $1 offset $2;`

const qUrlsUnfetched = `
select
  url, created, updated, last_head, last_get, status, content_type, content_sniff,
  content_length, file_name, title, id, headers_took, download_took, headers, meta, hash
from urls
where 
  last_get is null 
order by created desc 
limit $1 offset $2;`

const qUrlsForHash = `
select
  url, created, updated, last_head, last_get, status, content_type, content_sniff,
  content_length, file_name, title, id, headers_took, download_took, headers, meta, hash
from urls
where 
  hash = $1;`

const qUrlExistsByUrlString = `SELECT exists(SELECT 1 FROM urls WHERE url = $1)`

const qUrlByUrlString = `
select
  url, created, updated, last_head, last_get, status, content_type, content_sniff,
  content_length, file_name, title, id, headers_took, download_took, headers, meta, hash
from urls 
where
  url = $1;`

const qUrlById = `
select
  url, created, updated, last_head, last_get, status, content_type, content_sniff,
  content_length, file_name, title, id, headers_took, download_took, headers, meta, hash
from urls 
where
  id = $1;`

const qUrlExistsById = `SELECT exists(SELECT 1 FROM urls WHERE id = $1)`

const qUrlByHash = `
select
  url, created, updated, last_head, last_get, status, content_type, content_sniff,
  content_length, file_name, title, id, headers_took, download_took, headers, meta, hash
from urls 
where
  hash = $1;`

const qUrlExistsByHash = `SELECT exists(SELECT 1 FROM urls WHERE hash = $1)`

const qUrlInsert = `
insert into urls
  (url, created, updated, last_head, last_get, status, content_type, content_sniff,
  content_length, file_name, title, id, headers_took, download_took, headers, meta, hash)
values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17);`

const qUrlUpdate = `
update urls 
set 
  created=$2, updated=$3, last_head=$4, last_get=$5, status=$6, content_type=$7, content_sniff=$8,
  content_length=$9, file_name=$10, title=$11, id=$12, headers_took=$13, download_took=$14, headers=$15, meta=$16, hash=$17 
where 
  url = $1;`

const qUrlDelete = `
delete from urls 
where
  url = $1;`

const qUrlInboundLinkUrlStrings = `
select src 
from links
where
  dst = $1;`

const qUrlOutboundLinkUrlStrings = `
select dst 
from links 
where
  src = $1;`

const qUrlDstLinks = `
select 
  urls.url, urls.created, urls.updated, last_head, last_get, status, content_type, content_sniff, 
  content_length, file_name, title, id, headers_took, download_took, headers, meta, hash 
from urls, links
where 
  links.src = $1 and 
  links.dst = urls.url;`

// select all destination links that lead to content urls
const qUrlDstContentLinks = `
SELECT 
  urls.url, urls.created, urls.updated, last_head, last_get, status, content_type, content_sniff, 
  content_length, file_name, title, id, headers_took, download_took, headers, meta, hash 
FROM urls, links
WHERE 
  links.src = $1 AND 
  links.dst = urls.url AND
  urls.hash != '' AND
  urls.hash != '1220e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855' AND
  urls.content_sniff != 'text/html; charset=utf-8'
;`

const qUrlSrcLinks = `
select
  urls.url, urls.created, urls.updated, last_head, last_get, status, content_type, content_sniff, 
  content_length, file_name, title, id, headers_took, download_took, headers, meta, hash 
from urls, links 
where 
  links.dst = $1 and 
  links.src = urls.url;`

const qUncrawlableCreateTable = `
CREATE TABLE IF NOT EXISTS uncrawlables (
  id               text NOT NULL default '',
  url              text PRIMARY KEY NOT NULL,
  created          timestamp NOT NULL default (now() at time zone 'utc'),
  updated          timestamp NOT NULL default (now() at time zone 'utc'),
  creator_key_id   text NOT NULL default '',
  name             text NOT NULL default '',
  email            text NOT NULL default '',
  event_name       text NOT NULL default '',
  agency_name      text NOT NULL default '',
  agency_id        text NOT NULL default '',
  subagency_id     text NOT NULL default '',
  org_id           text NOT NULL default '',
  suborg_id        text NOT NULL default '',
  subprimer_id     text NOT NULL default '',
  ftp              boolean default false,
  database         boolean default false,
  interactive      boolean default false,
  many_files       boolean default false,
  comments         text NOT NULL default '',
  deleted          boolean NOT NULL default false
);`

const qUncrawlableExists = `SELECT exists(SELECT 1 FROM uncrawlables WHERE id = $1)`
const qUncrawlableExistsByUrl = `SELECT exists(SELECT 1 FROM uncrawlables WHERE url = $1)`

const qUncrawlablesList = `
select 
  id, url,created,updated,creator_key_id,
  name,email,event_name,agency_name,
  agency_id,subagency_id,org_id,suborg_id,subprimer_id,
  ftp,database,interactive,many_files,
  comments
from uncrawlables
order by created desc
limit $1 offset $2;`

const qUncrawlableInsert = `
insert into uncrawlables 
  ( id, url, created,updated,creator_key_id,
    name,email,event_name,agency_name,
    agency_id,subagency_id,org_id,suborg_id,subprimer_id,
    ftp,database,interactive,many_files,
    comments) 
values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19);`

const qUncrawlableUpdate = `
update uncrawlables 
set
  url = $2, created = $3, updated = $4, creator_key_id = $5,
  name = $6, email = $7, event_name = $8, agency_name = $9,
  agency_id = $10, subagency_id = $11, org_id = $12, suborg_id = $13, subprimer_id = $14,
  ftp = $15, database = $16,interactive = $17, many_files = $18,
  comments = $19
where id = $1;`

const qUncrawlableByUrl = `
SELECT 
  id, url,created,updated,creator_key_id,
  name,email,event_name,agency_name,
  agency_id,subagency_id,org_id,suborg_id,subprimer_id,
  ftp,database,interactive,many_files,
  comments
FROM uncrawlables 
WHERE url = $1;`

const qUncrawlableById = `
select 
  id, url,created,updated,creator_key_id,
  name,email,event_name,agency_name,
  agency_id,subagency_id,org_id,suborg_id,subprimer_id,
  ftp,database,interactive,many_files,
  comments
from uncrawlables 
where id = $1;`

const qUncrawlableDelete = `
delete from uncrawlables 
where url = $1;`

const qCustomCrawlCreateTable = `
CREATE TABLE IF NOT EXISTS custom_crawls (
  id               UUID PRIMARY KEY NOT NULL,
  created          timestamp NOT NULL default (now() at time zone 'utc'),
  updated          timestamp NOT NULL default (now() at time zone 'utc'),
  jwt              text NOT NULL default '',
  morphRunId       text NOT NULL default '',
  dateCompleted    timestamp NOT NULL default (now() at time zone 'utc'),
  githubRepo       text NOT NULL default '',
  originalUrl      text NOT NULL default '',
  sqliteChecksum   text NOT NULL default ''
);`

const qCustomCrawlExists = `SELECT exists(SELECT 1 FROM custom_crawls WHERE id = $1)`

const qCustomCrawlsList = `
select 
  id, created, updated,
  jwt, morphRunId, dateCompleted, githubRepo, originalUrl,
  sqliteChecksum
from custom_crawls
order by created DESC
limit $1 offset $2`

const qCustomCrawlInsert = `
insert into custom_crawls 
  (id, created, updated,
  jwt, morphRunId, dateCompleted, githubRepo, originalUrl,
  sqliteChecksum) 
values ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

const qCustomCrawlUpdate = `
update custom_crawls 
set
  created = $2, updated = $3,
  jwt = $4, morphRunId = $5, dateCompleted = $6, githubRepo = $7, originalUrl = $8,
  sqliteChecksum = $9
where id = $1`

const qCustomCrawlByUrl = `
SELECT 
  id, created, updated,
  jwt, morphRunId, dateCompleted, githubRepo, originalUrl,
  sqliteChecksum
FROM custom_crawls 
WHERE url = $1;`

const qCustomCrawlById = `
select 
  id, created, updated,
  jwt, morphRunId, dateCompleted, githubRepo, originalUrl,
  sqliteChecksum
from custom_crawls 
where id = $1;`

const qCustomCrawlDelete = `
delete from custom_crawls where id = $1;`
