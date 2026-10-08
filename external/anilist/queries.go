package anilist

const browseAnimeQuery = `
query ($page: Int, $perPage: Int) {
  Page(page: $page, perPage: $perPage) {
    pageInfo {
      currentPage
      hasNextPage
      lastPage
      perPage
      total
    }

    media(
      type: ANIME
      sort: POPULARITY_DESC
    ) {
      id
      title {
        romaji
        english
        native
      }
      description
      coverImage {
        medium
        large
      }
      episodes
      status
      format
      season
      seasonYear
      genres
      averageScore
    }
  }
}
`

const searchAnimeQuery = `
query (
  $page: Int
  $perPage: Int
  $id: Int
  $search: String
  $season: MediaSeason
  $seasonYear: Int
  $genre: String
  $status: MediaStatus
  $format: MediaFormat
) {
  Page(page: $page, perPage: $perPage) {
    pageInfo {
      currentPage
      hasNextPage
      lastPage
      perPage
      total
    }

    media(
      type: ANIME
      id: $id
      search: $search
      season: $season
      seasonYear: $seasonYear
      genre: $genre
      status: $status
      format: $format
      sort: POPULARITY_DESC
    ) {
      id
      title {
        romaji
        english
        native
      }
      description
      coverImage {
        medium
        large
      }
      episodes
      status
      format
      season
      seasonYear
      genres
      averageScore
    }
  }
}
`

const animeDetailQuery = `
query ($id: Int!) {
  Media(id: $id, type: ANIME) {
    id
    title {
      romaji
      english
      native
    }
    description
    coverImage {
      medium
      large
    }
    episodes
    status
    format
    season
    seasonYear
    genres
    averageScore
  }
}
`
