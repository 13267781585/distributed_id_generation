namespace go uuid.generator.server

struct Base {
    1: optional i64 code,
    2: optional string message,
}

struct GetUUIDBoundsRequest {
    1: optional i64 biz_code,
    2: optional i64 count,
}

struct UUIDBound {
    1: optional i64 start,
    2: optional i64 end,
}

struct GetUUIDBoundsResponse {
    1: optional list<UUIDBound> uuid_bounds,
    255: optional Base base,
}

service UUIDGeneratorServer {
    GetUUIDBoundsResponse GetUUIDBounds(1: GetUUIDBoundsRequest req)
}