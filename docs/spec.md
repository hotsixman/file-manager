# 파일 매니저

파일 매니저는 파일 시스템 리소스의 동시성 제어를 위해 파일 및 폴더의 상태를 관리하고 잠금(Lock) 기능을 제공한다.

## 잠금 상태 (Lock Status)

파일 및 폴더는 다음 중 하나의 상태를 가진다.

- **0 Idle**: 유휴 상태. 잠겨 있지 않음.
- **1 Reading**: 읽기 중. (다운로드 포함)
- **2 Moving From**: 이동/복사 소스. (이름 변경 포함)
- **3 Moving To**: 이동/복사 대상. (이름 변경, 업로드 포함)
- **4 Deleting**: 삭제 중.
- **5 Locked**: 명시적으로 잠긴 상태.

## 상태 코드 (Status Codes)

API 응답의 `status` 메타데이터 필드에 포함되는 코드들이다.

- **200 OK**: 요청 성공.
- **204 No Content**: 이미 잠금 해제된 상태.
- **400 Bad Request**: 잘못된 상태 값 입력 또는 절대 경로가 아님.
- **422 Unprocessable Entity**: 조상(Ancestor) 요소가 이미 잠겨 있어 현재 요소의 잠금이 불가능함.
- **429 Too Many Requests**: 이미 다른 상태로 잠겨 있음.
- **500 Internal Server Error**: 서버 내부 오류.

## API 정의

모든 통신은 `ndj-flow-client-go`를 통한 UDS(Unix Domain Socket)를 사용한다.

### 1. 기본 제어

#### lock
특정 리소스를 지정된 상태로 잠근다.
- **Command**: `lock`
- **Request Body**:
  ```json
  {
      "path": string,
      "status": 1 | 2 | 3 | 4 | 5
  }
  ```
- **Response Body**:
  ```json
  [
      {
          "status": number
      }
  ]
  ```

#### unlock
특정 리소스의 잠금을 해제한다.
- **Command**: `unlock`
- **Request Body**:
  ```json
  {
      "path": string
  }
  ```
- **Response Body**: 없음

#### check
특정 리소스의 현재 잠금 상태를 확인한다.
- **Command**: `check`
- **Request Body**:
  ```json
  {
      "path": string
  }
  ```
- **Response Body**:
  ```json
  [
      {
          "status": number,
          "blocked": boolean,
          "ancestorLocked": boolean,
          "decendentLocked": boolean
      }
  ]
  ```

### 2. 복합 작업 (Pre/Post)

작업 전(Pre)에 필요한 리소스를 잠그고, 작업 후(Post)에 잠금을 해제하는 API들이다.

#### 이동 / 복사 / 이름 변경
- **Commands**: `movePre`, `movePost`, `copyPre`, `copyPost`, `renamePre`, `renamePost`
- **Request Body**:
  ```json
  {
      "src": string,
      "dest": string
  }
  ```
- **Pre 동작**: `src`를 **2 (Moving From)**로, `dest`를 **3 (Moving To)**으로 잠금.
- **Post 동작**: `src`와 `dest` 모두 잠금 해제.

#### 삭제 / 읽기 / 업로드 / 다운로드
- **Commands**: `deletePre`/`Post`, `readPre`/`Post`, `uploadPre`/`Post`, `downloadPre`/`Post`
- **Request Body**:
  ```json
  {
      "path": string
  }
  ```
- **Pre 동작**:
    - `deletePre`: **4 (Deleting)** 상태로 잠금.
    - `readPre`, `downloadPre`: **1 (Reading)** 상태로 잠금.
    - `uploadPre`: **3 (Moving To)** 상태로 잠금.
- **Post 동작**: 해당 경로의 잠금 해제.

## 주의 사항

- 모든 경로는 절대 경로여야 한다.
- 내부적으로 `sync.Mutex`를 사용하여 스레드 안전성을 보장한다.
- 조상 요소가 잠겨 있는 경우 하위 요소의 잠금은 거부된다 (422).
- 자손 요소가 잠겨 있는 상태에서 조상 요소를 잠그는 것은 허용되나, 이 경우 자손의 잠금이 풀려도 조상의 잠금은 유지된다.
