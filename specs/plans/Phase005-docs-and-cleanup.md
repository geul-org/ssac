✅ 완료

# Phase 005: 문서 업데이트 + 정리

## 목표

1. 문서를 v2 문법에 맞게 전면 업데이트
2. CLAUDE.md를 v2 기준으로 재작성
3. v1 코드 정리 (v1/ 폴더 유지 여부 결정)

## 변경 사항

### 1. artifacts/manual-for-ai.md — 전면 재작성

- DSL 문법 섹션: v2 한 줄 표현식 문법
- 시퀀스 타입: 10개 (get, post, put, delete, empty, exists, state, auth, call, response)
- 코드젠 기능: gin 프레임워크 기준, @response 필드 매핑, @state/@auth JSON 입력
- 디렉토리 구조: v1/ 참조 제거

### 2. artifacts/manual-for-human.md — 전면 재작성

- 전체 문법 레퍼런스 v2 기준
- 코드젠 예시 전면 교체
- v1 → v2 마이그레이션 가이드 포함

### 3. README.md — 전면 재작성

- Core Idea 예시: v2 문법 + gin 코드젠 결과
- Sequence Types 테이블: v2 기준
- Code Generation Features: @response 필드 매핑, @state/@auth JSON 입력
- 프로젝트 구조: v2 기준

### 4. CLAUDE.md — v2 기준 재작성

- DSL 문법 섹션: v2 한 줄 표현식
- 타입별 필수 태그 테이블: v2 기준
- 코드젠 기능: v2 기준
- 디렉토리 구조: v2 기준

### 5. v1/ 폴더 보존

v1/ 폴더는 삭제하지 않고 그대로 유지한다. (참조용)

## 수정 파일

| 파일 | 내용 |
|---|---|
| `artifacts/manual-for-ai.md` | v2 기준 전면 재작성 |
| `artifacts/manual-for-human.md` | v2 기준 전면 재작성 |
| `README.md` | v2 기준 전면 재작성 |
| `CLAUDE.md` | v2 기준 재작성 |

## 의존성

- Phase 001-004 완료 후

## 검증

- 문서 내 코드 예시가 실제 파서/코드젠과 일치하는지 확인
- `ssac parse`/`ssac validate`/`ssac gen` 실행 결과와 문서 비교
