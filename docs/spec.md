# 파일 매니저
파일 매니저는 다음 기능의 구현을 목표로 한다.
- 파일을 사용하지 못하도록 잠굼

## 잠금
파일 잠금은 다음 상태를 가진다.
- `0` Idle
- `1` Reading
- `2` Moving From
    - 파일 복사, 이름 변경도 포함
- `3` Moving To
    - 파일 복사, 이름 변경도 포함
- `4` Deleting
- `5` Locked

폴더 잠금은 다음 상태를 가진다.
- `0` Idle
- `1` Listing
- `2` Moving From
    - 파일 복사, 이름 변경도 포함
- `3` Moving To
    - 파일 복사, 이름 변경도 포함
- `4` Deleting
- `5` Locked

## API
### `/lock`
```ts
{
    request: {
        param: {
            path: string,
            status: 1 | 2 | 3 | 4 | 5
        }
    },
    response: {
        body: {
            status: 0 | 1 | 2 | 3 | 4 | 5
        },
        status: 200 | 422 | 429
        // 200: 잠금 성공
        // 422: 하위,상위 요소가 잠겨 있어 잠금이 불가능함
        // 429: 이미 잠김
    }
}
```
특정 리소스를 잠금

### `/unlock`
```ts
{
    request: {
        param: {
            path: string
        }
    },
    response: {
        body: 0 | 1 | 2 | 3 | 4 | 5,
        status: 200 | 204 | 422
        // 200: 잠금 해제 성공
        // 204: 이미 잠금 해제됨
    }
}
```
특정 리소스의 잠금 해제

### `check`
```ts
{
    request: {
        param: {
            path: string
        }
    },
    response: {
        body: {
            status: number,
            blocked: boolean,
            ancestorLocked: boolean,
            decendentLocked: boolean
        }
    }
}
```

## 주의 사항
mutex를 사용할 것