# 📚 Danh mục Agents & Skills mặc định

> 🇬🇧 English version: [CATALOG.md](CATALOG.md)

> **Phạm vi:** Tài liệu này mô tả chi tiết các bộ agent/skill được tích hợp sẵn: bộ **Game Studios** (`dotagent-game-studios-*`, mới) và catalog `da-*`/`ds-*` trước đây (VoltAgent + mattpocock). Danh sách đầy đủ 842 built-in skills từ 59 vendor, số lượng theo category và cách bật `dotagent-*` được duy trì trong [README](../README.md#built-in-skills).

> Tất cả agent/skill đều bị tắt theo mặc định — bạn chọn bật cái nào và cho nền tảng nào trong `config.yaml` hoặc qua Web Dashboard.

---

## 🎮 Game Studios (49 agents · 73 skills)

Bộ **Game Studios** được import từ [donchitos/claude-code-game-studios](https://github.com/donchitos/claude-code-game-studios) (MIT) — một "studio game ảo" đầy đủ vai trò: chuyên gia engine (Unity, Unreal, Godot), lập trình viên gameplay/engine/network/AI/UI, designer, producer, QA, release, live-ops, narrative, audio, art. Kèm theo là các quy trình sprint / gate / release theo kiểu studio thật.

- **Agent:** `dotagent-game-studios-<tên>` (category `game-development`)
- **Skill:** thư mục `dotagent-game-studios-<tên>/`, frontmatter `dotagent:game-studios:<tên>` (category `Game Development`, vendor `game-studios`)
- Import lại từ upstream bằng `python3 scripts/import-game-studios.py`.

> ⚠️ Nhiều skill trong bộ này giả định cấu trúc dự án của CCGS (`design/`, `production/`, `.claude/rules/*`, các template trong `.claude/docs/templates/*`). Hãy chạy `/start` hoặc `/adopt` trước để khởi tạo/kiểm tra scaffolding.

### 🎯 Ban lãnh đạo & định hướng (5 agents)

| Agent | Mô tả | Sử dụng khi nào |
|---|---|---|
| **`creative-director`** | Thẩm quyền sáng tạo cao nhất: vision, tone, thẩm mỹ, giải quyết xung đột thiết kế | Cần quyết định vision game, chốt hướng thẩm mỹ, phân xử tranh luận thiết kế |
| **`technical-director`** | Quyết định kỹ thuật cấp cao: kiến trúc engine, lựa chọn công nghệ, chiến lược hiệu năng, rủi ro kỹ thuật | Chọn engine/tech stack, lập kiến trúc tổng thể, đánh giá rủi ro kỹ thuật |
| **`lead-programmer`** | Kiến trúc mức code, coding standard, code review, phân việc cho programmer chuyên biệt | Cần review code, thiết kế API, chiến lược refactor, dịch design thành cấu trúc code |
| **`producer`** | Sprint planning, milestone, quản lý rủi ro, đàm phán scope, điều phối liên phòng ban | Lập sprint, theo dõi milestone, xử lý scope creep |
| **`art-director`** | Bản sắc hình ảnh: style guide, art bible, chuẩn asset, palette, pipeline sản xuất art | Cần art bible, chuẩn hóa asset, review UI/UX visual |

### 🎲 Thiết kế game (7 agents)

| Agent | Mô tả | Sử dụng khi nào |
|---|---|---|
| **`game-designer`** | Thiết kế cơ chế và hệ thống: core loop, progression, combat, economy, trải nghiệm người chơi | Cần thiết kế/đánh giá cơ chế gameplay, viết GDD |
| **`systems-designer`** | Thiết kế chi tiết từng subsystem: công thức combat, đường cong progression, crafting, status effect | Cần spec toán học/công thức cho một hệ thống cụ thể |
| **`economy-designer`** | Kinh tế tài nguyên, loot, progression curve, in-game market | Thiết kế loot table, faucet/sink, cân bằng kinh tế |
| **`level-designer`** | Thiết kế không gian, bố trí encounter, pacing, environmental storytelling | Cần layout level, kế hoạch pacing, hướng dẫn kể chuyện qua môi trường |
| **`live-ops-designer`** | Chiến lược nội dung hậu phát hành: sự kiện mùa, battle pass, cadence, retention | Lập kế hoạch live service, sự kiện, retention mechanics |
| **`ux-designer`** | UX flow, interaction design, accessibility, information architecture, input | Thiết kế user flow, HUD, wireframe, xử lý input |
| **`accessibility-specialist`** | Đảm bảo game chơi được với đông đảo người chơi nhất; chuẩn accessibility, review UI | Cần audit accessibility, tùy chọn hỗ trợ (colorblind, remap, subtitle…) |

### 💻 Lập trình (9 agents)

| Agent | Mô tả | Sử dụng khi nào |
|---|---|---|
| **`gameplay-programmer`** | Hiện thực cơ chế game, player system, combat, tính năng tương tác | Cần code cơ chế đã được thiết kế |
| **`engine-programmer`** | Hệ thống lõi engine: rendering, physics, memory, resource loading, scene management | Làm việc ở tầng engine/framework |
| **`ai-programmer`** | AI game: behavior tree, state machine, pathfinding, perception, hành vi NPC | Xây AI cho enemy/NPC, hệ thống ra quyết định |
| **`network-programmer`** | Multiplayer networking: state replication, lag compensation, matchmaking, protocol | Xây netcode, đồng bộ state, matchmaking |
| **`ui-programmer`** | Hệ thống UI: menu, HUD, inventory, dialogue box, UI framework | Hiện thực màn hình/HUD từ UX spec |
| **`tools-programmer`** | Công cụ nội bộ: editor extension, content authoring tool, debug utility, pipeline automation | Cần tool cho designer/artist, tự động hóa pipeline |
| **`technical-artist`** | Cầu nối art–engineering: shader, VFX, tối ưu rendering, art pipeline tool | Cần shader/VFX, tối ưu hình ảnh, tool cho artist |
| **`prototyper`** | Prototype nhanh (concept prototype sau brainstorm, feature prototype trước GDD) | Cần chứng minh ý tưởng vui trước khi đầu tư thiết kế đầy đủ |
| **`devops-engineer`** | Build pipeline, CI/CD, version control workflow, deployment | Cần build script, CI, quy trình branch |

### 🎮 Chuyên gia engine (15 agents)

| Agent | Mô tả | Sử dụng khi nào |
|---|---|---|
| **`unity-specialist`** | Thẩm quyền về pattern/API/tối ưu Unity; MonoBehaviour vs DOTS | Mọi câu hỏi Unity tổng quát |
| **`unity-dots-specialist`** | DOTS/ECS: Entity Component System, Jobs, Burst | Cần hiệu năng data-oriented trong Unity |
| **`unity-shader-specialist`** | Shader Graph, HLSL, VFX Graph, URP/HDRP | Tùy biến rendering trong Unity |
| **`unity-ui-specialist`** | UI Toolkit (UXML/USS), UGUI, data binding, hiệu năng UI runtime | Xây UI Unity |
| **`unity-addressables-specialist`** | Addressables: group, load/unload, memory, catalog, remote content | Quản lý asset/nội dung tải động trong Unity |
| **`unreal-specialist`** | Thẩm quyền về pattern/API/tối ưu Unreal; Blueprint vs C++ | Mọi câu hỏi Unreal tổng quát |
| **`ue-blueprint-specialist`** | Kiến trúc Blueprint, ranh giới Blueprint/C++, tối ưu graph | Giữ Blueprint dễ bảo trì, quyết định chuyển sang C++ |
| **`ue-gas-specialist`** | Gameplay Ability System: ability, effect, attribute set, tag, prediction | Hiện thực combat/ability bằng GAS |
| **`ue-replication-specialist`** | Networking Unreal: property replication, RPC, prediction, relevancy, bandwidth | Multiplayer trên Unreal |
| **`ue-umg-specialist`** | UMG/CommonUI: widget, data binding, input routing, styling | Xây UI Unreal |
| **`godot-specialist`** | Thẩm quyền Godot; GDScript vs C# vs GDExtension | Mọi câu hỏi Godot tổng quát |
| **`godot-gdscript-specialist`** | GDScript: static typing, pattern, signal, coroutine, tối ưu | Viết/refactor GDScript |
| **`godot-csharp-specialist`** | C# trong Godot 4: .NET pattern, export attribute, signal delegate, async | Dùng C# trong Godot |
| **`godot-gdextension-specialist`** | GDExtension: binding C/C++/Rust, tối ưu native | Cần code native cho Godot |
| **`godot-shader-specialist`** | Godot shading language, visual shader, material, particle, post-processing | Tùy biến rendering trong Godot |

### 📖 Narrative & nội dung (4 agents)

| Agent | Mô tả | Sử dụng khi nào |
|---|---|---|
| **`narrative-director`** | Kiến trúc câu chuyện, world-building, nhân vật, chiến lược dialogue | Lập story arc, phát triển nhân vật |
| **`world-builder`** | Lore chi tiết: phe phái, văn hóa, lịch sử, địa lý, sinh thái | Xây dựng thế giới và quy tắc của nó |
| **`writer`** | Dialogue, lore entry, mô tả item, text môi trường | Cần văn bản người chơi nhìn thấy |
| **`localization-lead`** | Kiến trúc i18n, string table, locale testing, pipeline dịch | Cần hệ thống đa ngôn ngữ, quy trình dịch |

### 🔊 Audio (2 agents)

| Agent | Mô tả | Sử dụng khi nào |
|---|---|---|
| **`audio-director`** | Bản sắc âm thanh: music direction, triết lý sound design, chiến lược implementation, mix | Định hướng âm nhạc/âm thanh tổng thể |
| **`sound-designer`** | Spec chi tiết cho SFX, audio event, tham số mix | Cần spec sheet SFX, danh sách audio event |

### 🧪 QA, hiệu năng, bảo mật, phát hành & cộng đồng (7 agents)

| Agent | Mô tả | Sử dụng khi nào |
|---|---|---|
| **`qa-lead`** | Chiến lược test, triage bug, quality gate, quy trình test | Lập test plan, đánh giá severity, gate release |
| **`qa-tester`** | Test case chi tiết, bug report, checklist | Viết test case, regression checklist, báo bug |
| **`performance-analyst`** | Profiling, tìm bottleneck, đề xuất tối ưu, theo dõi metric | Game chậm, cần budget hiệu năng |
| **`security-engineer`** | Chống cheat, exploit, rò rỉ dữ liệu; review lỗ hổng, bảo vệ save | Cần anti-cheat, bảo mật save/network |
| **`release-manager`** | Pipeline release: checklist chứng nhận, submit store, yêu cầu platform, versioning | Chuẩn bị phát hành, submit lên store |
| **`analytics-engineer`** | Telemetry, tracking hành vi người chơi, A/B test, data pipeline | Cần đo lường/phân tích người chơi |
| **`community-manager`** | Giao tiếp với người chơi: patch note, social, thu thập feedback, triage bug từ cộng đồng | Viết patch note, quản lý phản hồi cộng đồng |

### 🛠 Skills — Khởi động & khám phá (6 skills)

| Skill | Mô tả | Sử dụng khi nào |
|---|---|---|
| **`start`** | Onboarding lần đầu — hỏi bạn đang ở đâu rồi dẫn tới workflow phù hợp | Mới bắt đầu, chưa biết làm gì |
| **`help`** | Phân tích việc đã làm và câu hỏi của bạn, gợi ý bước kế tiếp | "Tiếp theo nên làm gì?", đang bí |
| **`adopt`** | Onboarding dự án brownfield: audit artifact hiện có theo template, phân loại thiếu sót | Đưa dự án đang có vào quy trình CCGS |
| **`onboard`** | Sinh tài liệu onboarding cho contributor/agent mới | Có người/agent mới tham gia dự án |
| **`project-stage-detect`** | Tự phát hiện giai đoạn dự án, thiếu sót, và gợi ý bước tiếp | "Dự án đang ở đâu?" |
| **`setup-engine`** | Cấu hình engine + phiên bản, ghim vào CLAUDE.md, bổ sung tài liệu tham chiếu engine | Ngay sau brainstorm, trước prototype |

### 🎲 Skills — Concept & thiết kế (12 skills)

| Skill | Mô tả | Sử dụng khi nào |
|---|---|---|
| **`brainstorm`** | Ý tưởng game có hướng dẫn: từ con số 0 đến concept document có cấu trúc | Bắt đầu một game mới |
| **`prototype`** | Concept prototype để xác nhận ý tưởng đáng thiết kế trước khi viết GDD | Ngay sau `/brainstorm` + `/setup-engine` |
| **`map-systems`** | Phân rã concept thành các hệ thống, map dependency, ưu tiên thứ tự thiết kế, tạo systems index | Sau khi có concept, trước khi viết GDD |
| **`design-system`** | Viết GDD cho một hệ thống theo từng section có hướng dẫn | Cần GDD đầy đủ cho một hệ thống |
| **`quick-design`** | Spec thiết kế nhẹ cho thay đổi nhỏ (tuning, cơ chế phụ, balance) | GDD hệ thống đã có, chỉ cần điều chỉnh |
| **`design-review`** | Review GDD về tính đầy đủ, nhất quán, khả thi | Trước khi giao GDD cho lập trình |
| **`review-all-gdds`** | Review chéo toàn bộ GDD: mâu thuẫn, tham chiếu cũ | Trước milestone thiết kế / pre-production |
| **`consistency-check`** | Quét GDD so với entity registry để phát hiện lệch số liệu giữa tài liệu | Sau khi sửa nhiều GDD |
| **`propagate-design-change`** | Khi GDD thay đổi, tìm ADR/traceability bị ảnh hưởng, tạo change report | Vừa sửa GDD đã được kiến trúc hóa |
| **`balance-check`** | Phân tích file data/công thức để tìm outlier, progression hỏng, chiến lược thoái hóa | Sau khi sửa balance data |
| **`ux-design`** | Viết UX spec cho screen/flow/HUD theo section | Cần spec UX trước khi làm UI |
| **`ux-review`** | Kiểm tra UX spec/HUD về đầy đủ, accessibility, khớp GDD, sẵn sàng hiện thực | Trước khi giao UX spec cho UI programmer |

### 🎨 Skills — Art & asset (3 skills)

| Skill | Mô tả | Sử dụng khi nào |
|---|---|---|
| **`art-bible`** | Viết Art Bible theo section — spec bản sắc hình ảnh gate mọi sản xuất asset | Sau brainstorm, trước sản xuất art |
| **`asset-spec`** | Sinh spec hình ảnh + prompt AI generation cho từng asset từ GDD/level doc/character profile | Cần brief asset cho artist/AI |
| **`asset-audit`** | Audit asset về naming, size budget, format, pipeline; tìm asset mồ côi | Trước milestone, sau import asset lớn |

### 🏛 Skills — Kiến trúc & lập kế hoạch (7 skills)

| Skill | Mô tả | Sử dụng khi nào |
|---|---|---|
| **`create-architecture`** | Viết master architecture document theo section, đọc toàn bộ GDD/ADR/engine | Sau khi GDD được duyệt |
| **`architecture-decision`** | Tạo ADR: bối cảnh, phương án, hệ quả | Có quyết định kỹ thuật quan trọng |
| **`architecture-review`** | Kiểm tra kiến trúc so với GDD, dựng traceability matrix | Sau create-architecture, trước create-epics |
| **`create-control-manifest`** | Sinh bảng quy tắc phẳng cho programmer: phải làm / không được làm, theo hệ thống | Sau khi kiến trúc hoàn tất |
| **`create-epics`** | Dịch GDD + kiến trúc thành epic (mỗi module một epic) | Bắt đầu production planning |
| **`create-stories`** | Chia một epic thành story file có thể hiện thực, nhúng requirement GDD/ADR | Trước sprint |
| **`reverse-document`** | Sinh tài liệu design/kiến trúc từ code có sẵn | Dự án thiếu tài liệu |

### 🏃 Skills — Sprint & production (11 skills)

| Skill | Mô tả | Sử dụng khi nào |
|---|---|---|
| **`sprint-plan`** | Tạo/cập nhật sprint plan dựa trên milestone, việc đã xong, capacity | Đầu sprint |
| **`sprint-status`** | Snapshot tiến độ sprint nhanh, đánh giá burndown | Giữa sprint |
| **`story-readiness`** | Kiểm tra story đã đủ điều kiện hiện thực (GDD, ADR, engine note, acceptance criteria) | Trước khi bắt đầu story |
| **`dev-story`** | Đọc story và hiện thực: nạp đầy đủ context, route tới programmer phù hợp | Làm một story |
| **`story-done`** | Review hoàn tất story: kiểm từng acceptance criterion, lệch GDD/ADR | Kết thúc story |
| **`estimate`** | Ước lượng effort theo độ phức tạp, dependency, velocity, rủi ro | Cần estimate có độ tin cậy |
| **`scope-check`** | Phát hiện scope creep so với kế hoạch gốc | Nghi ngờ phình scope |
| **`gate-check`** | Kiểm tra điều kiện chuyển giai đoạn → PASS/CONCERNS/FAIL | Trước khi sang phase mới |
| **`milestone-review`** | Review milestone: hoàn thiện, chất lượng, rủi ro, go/no-go | Cuối milestone |
| **`retrospective`** | Retro sprint/milestone: velocity, blocker, pattern | Cuối sprint |
| **`vertical-slice`** | Xây build end-to-end chất lượng production để xác nhận game loop | Pre-production |

### 🧑‍💻 Skills — Code & kỹ thuật (6 skills)

| Skill | Mô tả | Sử dụng khi nào |
|---|---|---|
| **`code-review`** | Review kiến trúc & chất lượng cho file/tập file | Trước khi merge |
| **`tech-debt`** | Theo dõi, phân loại, ưu tiên nợ kỹ thuật; duy trì debt register | Định kỳ / trước milestone |
| **`perf-profile`** | Quy trình profiling có cấu trúc, so với budget, đề xuất tối ưu | Game không đạt hiệu năng mục tiêu |
| **`security-audit`** | Audit lỗ hổng: save tampering, cheat, network exploit, lộ dữ liệu | Trước release |
| **`localize`** | Pipeline localization đầy đủ: quét string cứng, string table, validate, brief dịch giả | Chuẩn bị đa ngôn ngữ |
| **`content-audit`** | So sánh số lượng nội dung GDD với nội dung đã hiện thực | Theo dõi tiến độ nội dung |

### 🧪 Skills — Testing & QA (13 skills)

| Skill | Mô tả | Sử dụng khi nào |
|---|---|---|
| **`test-setup`** | Scaffold test framework + CI cho engine của dự án | Bắt đầu dự án |
| **`test-helpers`** | Sinh thư viện helper test đặc thù engine | Sau test-setup |
| **`qa-plan`** | Sinh QA test plan cho sprint/feature | Đầu sprint hoặc feature lớn |
| **`smoke-check`** | Gate smoke test critical path trước khi giao QA | Trước QA hand-off |
| **`regression-suite`** | Map coverage với critical path GDD, tìm bug đã sửa thiếu regression test | Duy trì bộ regression |
| **`test-flakiness`** | Phát hiện test không ổn định từ log CI | CI đỏ ngẫu nhiên |
| **`test-evidence-review`** | Review chất lượng test file và evidence thủ công | Trước gate/release |
| **`soak-test`** | Sinh protocol soak test cho phiên chơi dài | Tìm leak/degradation chậm |
| **`playtest-report`** | Template/phân tích báo cáo playtest | Sau playtest |
| **`bug-report`** | Tạo bug report có cấu trúc, đủ bước tái hiện | Gặp bug |
| **`bug-triage`** | Đọc toàn bộ bug, đánh giá lại priority/severity, gán sprint | Định kỳ triage |
| **`skill-test`** / **`skill-improve`** | Kiểm tra & cải thiện chính các skill CCGS (linter, spec, audit) | Đang chỉnh sửa skill |

### 🚀 Skills — Release & live-ops (6 skills)

| Skill | Mô tả | Sử dụng khi nào |
|---|---|---|
| **`release-checklist`** | Checklist pre-release: build, certification, store metadata | Chuẩn bị release |
| **`launch-checklist`** | Kiểm tra sẵn sàng launch mọi phòng ban + go/no-go | Trước ngày launch |
| **`hotfix`** | Quy trình sửa khẩn cấp có audit trail, bỏ qua sprint | Lỗi nghiêm trọng trên production |
| **`day-one-patch`** | Chuẩn bị patch ngày đầu: scope, ưu tiên, hiện thực, QA gate | Sau gold, trước launch |
| **`changelog`** | Sinh changelog nội bộ + người chơi từ git/sprint/design | Cuối sprint/release |
| **`patch-notes`** | Patch note cho người chơi từ git history/changelog | Mỗi bản cập nhật |

### 👥 Skills — Điều phối team (9 skills)

Mỗi skill `team-*` điều phối một nhóm agent làm việc cùng nhau theo pipeline hoàn chỉnh.

| Skill | Nhóm agent | Sử dụng khi nào |
|---|---|---|
| **`team-combat`** | game-designer, gameplay-programmer, ai-programmer, technical-artist, sound-designer, qa-tester | Thiết kế → hiện thực → test một hệ thống combat |
| **`team-level`** | level-designer, narrative-director, world-builder, art-director, systems-designer, qa-tester | Xây một area/level hoàn chỉnh |
| **`team-narrative`** | narrative-director, writer, world-builder, level-designer | Nội dung câu chuyện & lore |
| **`team-ui`** | ux-designer, ui-programmer, art-director, qa-tester (tích hợp `/ux-design`, `/ux-review`) | UX spec → visual → hiện thực → polish |
| **`team-audio`** | audio-director, sound-designer, technical-artist, gameplay-programmer | Pipeline audio từ định hướng đến implementation |
| **`team-polish`** | performance-analyst, technical-artist, sound-designer, qa-tester | Tối ưu & polish một feature |
| **`team-qa`** | qa-lead, qa-tester | Chu trình test đầy đủ |
| **`team-release`** | release-manager, qa-lead, devops-engineer, producer | Từ release candidate đến deploy |
| **`team-live-ops`** | live-ops-designer, economy-designer, analytics-engineer, community-manager | Kế hoạch nội dung hậu phát hành |


---

## 🛠 Skills (Slash Commands) — bộ mattpocock

Skills là các **quy trình làm việc có cấu trúc** (slash command) mà agent sẽ tuân theo khi được kích hoạt. Khác với agent (định nghĩa *agent là ai*), skill định nghĩa *agent làm gì theo quy trình nào*.

Mỗi skill được lưu dưới dạng thư mục `dotagent-mattpocock-<tên>/SKILL.md` (tên `ds-*` bên dưới là tên cũ, hiện tương ứng `dotagent-mattpocock-*`), có thể kèm thư mục `references/` chứa các file tham chiếu bổ sung.

### 🔧 Engineering (9 skills)

| Skill | Mô tả | Sử dụng khi nào |
|---|---|---|
| **`ds-diagnose`** | Vòng lặp chẩn đoán: Tái tạo → Thu nhỏ → Giả thuyết → Đo lường → Sửa → Test regression | Gặp bug khó, regression hiệu năng, user nói "debug này" hoặc "sửa lỗi này" |
| **`ds-grill-with-docs`** | Thử thách kế hoạch dựa trên domain model, làm sắc nét thuật ngữ, cập nhật CONTEXT.md và ADRs inline | Muốn stress-test kế hoạch dựa trên tài liệu và ngôn ngữ domain hiện có |
| **`ds-improve-codebase-architecture`** | Tìm cơ hội cải thiện kiến trúc dựa trên CONTEXT.md và docs/adr/ | Muốn refactor, gộp module chặt, làm codebase dễ test và dễ AI điều hướng |
| **`ds-setup-matt-pocock-skills`** | Thiết lập `## Agent skills` trong AGENTS.md và `docs/agents/` cho repo | **Chạy đầu tiên** — trước khi dùng `to-issues`, `triage`, `diagnose`, `tdd`, v.v. |
| **`ds-tdd`** | Phát triển hướng kiểm thử với vòng lặp red-green-refactor | Muốn xây tính năng hoặc sửa bug theo TDD, cần integration test |
| **`ds-to-issues`** | Chia kế hoạch/spec/PRD thành issue độc lập theo "tracer-bullet vertical slices" | Muốn chuyển kế hoạch thành ticket công việc trên issue tracker |
| **`ds-to-prd`** | Chuyển ngữ cảnh hội thoại hiện tại thành PRD | Muốn tạo PRD từ cuộc trò chuyện đang diễn ra |
| **`ds-triage`** | Phân loại issue qua state machine với các vai trò triage | Muốn tạo issue, phân loại bug/feature request, chuẩn bị issue cho agent tự động |
| **`ds-zoom-out`** | Yêu cầu agent cung cấp ngữ cảnh rộng hơn, góc nhìn cấp cao | Chưa quen với phần code, cần hiểu nó nằm ở đâu trong bức tranh tổng thể |

### 🎯 Productivity (3 skills)

| Skill | Mô tả | Sử dụng khi nào |
|---|---|---|
| **`ds-caveman`** | Đơn giản hóa giải thích xuống mức cơ bản nhất — kiểu "người hang động" | Cần hiểu nhanh một khái niệm phức tạp bằng ngôn ngữ đơn giản |
| **`ds-grill-me`** | Phỏng vấn user không ngừng, giải quyết từng nhánh cây quyết định | Muốn stress-test kế hoạch/thiết kế, hoặc nói "grill me" |
| **`ds-write-a-skill`** | Tạo skill mới với cấu trúc đúng, progressive disclosure, tài nguyên đi kèm | Muốn viết hoặc xây dựng một skill mới cho agent |

### 🔩 Misc (4 skills)

| Skill | Mô tả | Sử dụng khi nào |
|---|---|---|
| **`ds-git-guardrails`** | Thiết lập hooks chặn lệnh git nguy hiểm (push, reset --hard, clean, branch -D) | Muốn bảo vệ repo khỏi thao tác git phá hoại trong Claude Code |
| **`ds-migrate-to-shoehorn`** | Di chuyển file test từ `as` type assertion sang @total-typescript/shoehorn | Dự án TypeScript cần thay thế `as` trong test bằng partial test data |
| **`ds-scaffold-exercises`** | Tạo cấu trúc thư mục bài tập: sections, problems, solutions, explainers | Cần scaffold bài tập cho khóa học hoặc section mới |
| **`ds-setup-pre-commit`** | Thiết lập Husky pre-commit hooks với lint-staged, type checking, tests | Muốn thêm pre-commit hooks, Husky, hoặc kiểm tra tự động khi commit |

---

## 🤖 Agents — bộ VoltAgent

Agents là các **chuyên gia ảo** với vai trò và chuyên môn cụ thể. Mỗi agent chứa system prompt hướng dẫn AI hành xử theo đúng vai trò khi được kích hoạt.

Mỗi agent được lưu dưới dạng file `dotagent-voltagent-<tên>.md` với frontmatter YAML (description, category) và nội dung Markdown. Tên trong bảng là tên rút gọn.

### 💼 Business & Product (12 agents)

Các agent hỗ trợ nghiệp vụ kinh doanh, quản lý sản phẩm, và chiến lược.

| Agent | Mô tả | Sử dụng khi nào |
|---|---|---|
| **`business-analyst`** | Phân tích quy trình nghiệp vụ, thu thập yêu cầu | Cần phân tích yêu cầu, cải tiến quy trình, hoặc tạo tài liệu BRD |
| **`content-marketer`** | Chiến lược nội dung, marketing SEO, chiến dịch đa kênh | Cần lập kế hoạch nội dung, viết bài SEO, hoặc chạy chiến dịch marketing |
| **`customer-success-manager`** | Đánh giá sức khỏe khách hàng, chiến lược giữ chân | Cần phân tích churn, tối ưu retention, hoặc tìm cơ hội upsell |
| **`legal-advisor`** | Soạn hợp đồng, rà soát tuân thủ, bảo vệ IP | Cần review hợp đồng, đánh giá rủi ro pháp lý, hoặc tư vấn tuân thủ |
| **`license-engineer`** | Lựa chọn giấy phép OSI, tuân thủ dependency | Cần chọn license cho dự án, kiểm tra tuân thủ dependency |
| **`marketing-analyst`** | Hiệu suất chiến dịch, mô hình attribution | Cần phân tích ROI chiến dịch, xây dựng mô hình tăng trưởng |
| **`product-manager`** | Ưu tiên tính năng, lập roadmap, điều phối stakeholder | Cần lập kế hoạch sản phẩm, ưu tiên backlog, viết user story |
| **`sales-engineer`** | Pre-sales kỹ thuật, kiến trúc giải pháp, PoC | Cần demo kỹ thuật, thiết kế giải pháp cho khách, hoặc làm PoC |
| **`scrum-master`** | Sprint planning, retrospective, loại bỏ impediment | Cần tổ chức ceremony Agile, cải thiện velocity, hoặc gỡ impediment |
| **`technical-writer`** | Tài liệu API, hướng dẫn sử dụng, tài liệu SDK | Cần viết/cải thiện tài liệu kỹ thuật, API docs, hoặc getting-started |
| **`ux-researcher`** | Nghiên cứu người dùng, usability testing, persona | Cần khảo sát UX, phân tích hành vi user, hoặc validate thiết kế |
| **`wordpress-master`** | Kiến trúc WordPress, WooCommerce, bảo mật | Cần xây/tối ưu WordPress, custom theme/plugin, hoặc hardening |

### ⚙️ Core Development (11 agents)

Các agent cho phát triển phần mềm cốt lõi — backend, frontend, API, database.

| Agent | Mô tả | Sử dụng khi nào |
|---|---|---|
| **`api-designer`** | Đặc tả API, RESTful patterns, GraphQL schema | Cần thiết kế API mới, review contract, hoặc viết OpenAPI spec |
| **`backend-developer`** | API phía server, microservices, backend vững chắc | Cần xây dựng backend, thiết kế service, hoặc tối ưu server-side |
| **`database-architect`** | Thiết kế schema, tối ưu truy vấn, migration | Cần thiết kế DB, viết migration, hoặc tối ưu query chậm |
| **`frontend-developer`** | Frontend React/Vue/Angular, responsive design | Cần xây giao diện, xử lý state management, hoặc tối ưu UX |
| **`fullstack-developer`** | Phát triển ứng dụng end-to-end | Cần làm việc cả frontend lẫn backend cùng lúc |
| **`graphql-developer`** | GraphQL schema, resolver, federation | Cần thiết kế/triển khai GraphQL API, hoặc setup federation |
| **`legacy-modernizer`** | Hiện đại hóa hệ thống legacy | Cần migrate hệ thống cũ, strangler fig pattern, hoặc rewrite |
| **`low-level-designer`** | Thiết kế class-level OOP/functional, SOLID | Cần thiết kế class diagram, áp dụng design pattern |
| **`microservices-architect`** | Hệ thống phân tán, service mesh, event-driven | Cần tách monolith, thiết kế microservices, hoặc event sourcing |
| **`ui-designer`** | Giao diện trực quan, design system, component library | Cần thiết kế UI, xây design system, hoặc tạo component library |
| **`websocket-engineer`** | Giao tiếp realtime hai chiều ở quy mô lớn | Cần tính năng realtime: chat, notification, live update |

### 🧠 Data & AI (13 agents)

Các agent cho khoa học dữ liệu, machine learning, và trí tuệ nhân tạo.

| Agent | Mô tả | Sử dụng khi nào |
|---|---|---|
| **`ai-engineer`** | Hệ thống AI end-to-end, lựa chọn model, deployment | Cần xây pipeline AI, chọn model phù hợp, hoặc deploy model |
| **`computer-vision-engineer`** | Phân tích ảnh/video, object detection, OCR | Cần xử lý hình ảnh, nhận diện đối tượng, hoặc đọc text từ ảnh |
| **`data-engineer`** | Pipeline ETL, data warehouse, streaming | Cần xây data pipeline, thiết kế warehouse, hoặc streaming data |
| **`data-pipeline-architect`** | Hạ tầng dữ liệu quy mô lớn, xử lý realtime | Cần kiến trúc data platform, xử lý big data, hoặc realtime analytics |
| **`data-scientist`** | Mô hình thống kê, thí nghiệm ML, trực quan hóa | Cần phân tích dữ liệu, xây model dự đoán, hoặc A/B test |
| **`data-visualization`** | Dashboard tương tác, D3.js, Plotly | Cần tạo biểu đồ, dashboard, hoặc báo cáo trực quan |
| **`elasticsearch-specialist`** | Cluster tìm kiếm, tối ưu query, quản lý index | Cần thiết lập/tối ưu Elasticsearch, xây full-text search |
| **`etl-specialist`** | Pipeline trích xuất, biến đổi, nạp dữ liệu | Cần xây ETL job, transform data, hoặc sync giữa các hệ thống |
| **`llm-architect`** | Ứng dụng LLM, RAG, fine-tuning | Cần xây ứng dụng AI/chatbot, triển khai RAG, hoặc fine-tune model |
| **`ml-engineer`** | Phát triển model ML, training pipeline | Cần train model, xây ML pipeline, hoặc tối ưu inference |
| **`nlp-engineer`** | Xử lý văn bản, sentiment analysis | Cần phân tích ngôn ngữ, phân loại text, hoặc trích xuất entity |
| **`playwright-expert`** | Tự động hóa trình duyệt, E2E testing, scraping | Cần viết E2E test, scraping web, hoặc tự động hóa browser |
| **`prompt-engineer`** | Thiết kế prompt, chain-of-thought, evaluation | Cần viết/tối ưu prompt, xây evaluation framework cho LLM |

### 🔨 Developer Experience (14 agents)

Các agent cải thiện trải nghiệm và năng suất phát triển.

| Agent | Mô tả | Sử dụng khi nào |
|---|---|---|
| **`build-engineer`** | Hiệu suất build, tối ưu compilation, scaling | Build chậm, cần tối ưu CI build time hoặc scaling |
| **`cli-developer`** | Công cụ dòng lệnh và ứng dụng terminal | Cần xây CLI tool, TUI app, hoặc script phức tạp |
| **`documentation-engineer`** | Documentation-as-code, tài liệu kiến trúc | Cần hệ thống tài liệu tự động, docs site, hoặc arch docs |
| **`git-specialist`** | Git workflow nâng cao, branching, history | Cần giải quyết conflict phức tạp, rebase, hoặc thiết kế git flow |
| **`github-actions-specialist`** | CI/CD với GitHub Actions | Cần viết/tối ưu GitHub Actions workflow |
| **`ide-plugin-developer`** | Extension cho VS Code, JetBrains | Cần xây plugin IDE, code action, hoặc language server |
| **`json-wrangler`** | Biến đổi JSON/YAML, validation, jq | Cần transform data phức tạp, viết JSON schema, hoặc dùng jq |
| **`monorepo-engineer`** | Kiến trúc Nx/Turborepo/Lerna | Cần setup monorepo, quản lý workspace, hoặc tối ưu task |
| **`open-source-advisor`** | Đóng góp OSS, quản trị cộng đồng | Cần chiến lược OSS, viết CONTRIBUTING.md, quản lý contributor |
| **`refactoring-specialist`** | Refactoring code, giảm tech debt | Code chất lượng kém, cần restructure hoặc giảm complexity |
| **`regex-master`** | Pattern regex phức tạp, validation | Cần viết regex phức tạp, parse text, hoặc validate format |
| **`slack-expert`** | Ứng dụng Slack, bot, tích hợp API | Cần xây Slack bot, webhook, hoặc Slack app |
| **`tooling-engineer`** | Công cụ developer, code generator | Cần xây developer tool, scaffolding, hoặc code gen |
| **`vibe-coder`** | Prototyping nhanh, creative coding | Cần MVP nhanh, prototype, hoặc thử nghiệm ý tưởng |

### 🏗 Infrastructure (16 agents)

Các agent cho hạ tầng, DevOps, và quản trị hệ thống.

| Agent | Mô tả | Sử dụng khi nào |
|---|---|---|
| **`azure-infra-engineer`** | Hạ tầng Azure, networking, deployment | Cần triển khai/quản lý hạ tầng trên Azure |
| **`cicd-engineer`** | Pipeline CI/CD, tự động hóa deployment | Cần thiết kế pipeline CI/CD, tự động hóa release |
| **`cloud-architect`** | Thiết kế hạ tầng cloud, multi-cloud | Cần thiết kế cloud architecture, chọn dịch vụ cloud |
| **`devops-engineer`** | Tự động hóa hạ tầng, monitoring | Cần IaC, monitoring stack, hoặc deployment automation |
| **`docker-expert`** | Tối ưu container, multi-stage build | Cần viết Dockerfile, tối ưu image size, hoặc Compose |
| **`gcp-specialist`** | Dịch vụ và kiến trúc GCP | Cần triển khai trên Google Cloud |
| **`kubernetes-specialist`** | Quản lý cluster K8s, Helm, operator | Cần deploy K8s, viết Helm chart, hoặc custom operator |
| **`linux-sysadmin`** | Quản trị server Linux, shell scripting | Cần quản trị Linux server, viết bash script |
| **`network-engineer`** | Kiến trúc mạng, bảo mật, troubleshooting | Cần thiết kế network, firewall, hoặc debug connectivity |
| **`nginx-specialist`** | Cấu hình Nginx, load balancing, reverse proxy | Cần cấu hình Nginx, SSL, hoặc reverse proxy |
| **`powershell-admin`** | Tự động hóa Windows, Active Directory | Cần script PowerShell, quản lý AD, hoặc GPO |
| **`security-engineer`** | Zero-trust architecture, bảo mật CI/CD | Cần thiết kế security architecture, shift-left security |
| **`sre-engineer`** | SLO/SLI, error budget, chaos engineering | Cần định nghĩa SLO, giảm toil, hoặc chaos test |
| **`terraform-engineer`** | Infrastructure as code, Terraform | Cần viết/refactor Terraform, quản lý state |
| **`terragrunt-expert`** | Terragrunt orchestration, cấu hình DRY | Cần Terragrunt wrapper, multi-env deployment |
| **`windows-infra-admin`** | Windows Server, AD, Group Policy | Cần quản trị Windows Server, AD, DNS, DHCP |

### 💬 Language Specialists (30 agents)

Chuyên gia cho từng ngôn ngữ lập trình và framework cụ thể.

| Agent | Mô tả | Sử dụng khi nào |
|---|---|---|
| **`angular-architect`** | Angular 15+ enterprise | Xây ứng dụng Angular lớn, lazy loading, hoặc enterprise patterns |
| **`astro-developer`** | Astro framework, website content-driven | Xây website tĩnh/content-driven với Astro |
| **`cpp-systems-developer`** | C++ systems, quản lý bộ nhớ | Lập trình hệ thống C++, tối ưu performance, RAII |
| **`csharp-dotnet-developer`** | C#/.NET enterprise | Xây ứng dụng .NET, ASP.NET Core, hoặc Blazor |
| **`django-developer`** | Django web và REST API | Xây web Django, DRF API, hoặc admin site |
| **`elixir-phoenix-developer`** | Elixir/Phoenix realtime | Xây ứng dụng realtime, LiveView, hoặc distributed system |
| **`flutter-developer`** | Flutter cross-platform | Xây mobile/web app bằng Flutter, Dart |
| **`golang-pro`** | Go, concurrency, hiệu năng | Viết Go application, goroutine, hoặc tối ưu performance |
| **`java-enterprise-architect`** | Java enterprise, Spring | Xây hệ thống Java lớn, Spring Boot, microservices |
| **`kotlin-expert`** | Kotlin/Android, coroutine | Xây ứng dụng Android, Kotlin Multiplatform |
| **`laravel-expert`** | Laravel PHP, Eloquent ORM | Xây web Laravel, Livewire, hoặc API backend PHP |
| **`nestjs-architect`** | NestJS enterprise backend | Xây API NestJS, microservices TypeScript |
| **`nextjs-developer`** | Next.js full-stack, SSR/SSG | Xây web Next.js, App Router, Server Components |
| **`nuxt-specialist`** | Nuxt 3, hệ sinh thái Vue | Xây web Nuxt, auto-imports, server routes |
| **`perl-modernizer`** | Hiện đại hóa Perl, Moose | Cần modernize codebase Perl legacy |
| **`php-engineer`** | PHP 8+, framework, hiệu năng | Viết PHP hiện đại, Composer, hoặc tối ưu performance |
| **`python-pro`** | Python, async, xử lý dữ liệu | Viết Python app, FastAPI, hoặc data processing |
| **`r-statistician`** | R thống kê, phân tích dữ liệu | Phân tích thống kê, ggplot2, hoặc Shiny dashboard |
| **`rails-developer`** | Ruby on Rails | Xây web Rails, ActiveRecord, Hotwire |
| **`react-native-developer`** | React Native mobile | Xây mobile app React Native cross-platform |
| **`react-specialist`** | React 18+, hooks, state | Xây giao diện React, hooks phức tạp, state management |
| **`ruby-pro`** | Ruby, metaprogramming | Viết Ruby thuần, gem, hoặc DSL |
| **`rust-engineer`** | Rust systems, memory safety | Lập trình Rust, ownership, async, hoặc systems code |
| **`spring-boot-engineer`** | Spring Boot 3+ enterprise | Xây Spring Boot app, JPA, Security, Cloud |
| **`sql-pro`** | SQL, schema, indexing | Tối ưu query phức tạp, thiết kế schema, hoặc index |
| **`swift-expert`** | Swift/iOS/macOS native | Xây iOS/macOS app, SwiftUI, hoặc async/await |
| **`symfony-specialist`** | Symfony 6+/7+, Doctrine ORM | Xây ứng dụng Symfony, Messenger, API Platform |
| **`typescript-pro`** | TypeScript advanced types | Viết generic phức tạp, type-level programming |
| **`vue-expert`** | Vue 3, Composition API | Xây Vue app, Composition API, Pinia |
| **`wordpress-master`** | Theme, plugin, WooCommerce | Xây WordPress custom, Gutenberg blocks |

### 🎭 Meta-Orchestration (11 agents)

Các agent điều phối và quản lý agent khác.

| Agent | Mô tả | Sử dụng khi nào |
|---|---|---|
| **`agent-installer`** | Khám phá, duyệt, cài đặt agent | Cần tìm và cài đặt agent mới cho Claude Code |
| **`agent-organizer`** | Tổ chức và tối ưu nhóm multi-agent | Cần lập nhóm agent cho dự án phức tạp |
| **`codebase-orchestrator`** | Quản trị refactor toàn repo | Cần refactor lớn nhiều file với approval loop |
| **`context-manager`** | Quản lý trạng thái chia sẻ giữa agent | Nhiều agent cần đồng bộ dữ liệu với nhau |
| **`error-coordinator`** | Xử lý lỗi phối hợp giữa component | Lỗi phân tán xảy ra ở nhiều component cùng lúc |
| **`it-ops-orchestrator`** | Điều phối hoạt động IT đa lĩnh vực | Cần phối hợp nhiều domain: PS, AD, network, đồng thời |
| **`knowledge-synthesizer`** | Trích xuất pattern từ tương tác agent | Cần rút ra insight từ lịch sử agent hoạt động |
| **`multi-agent-coordinator`** | Phối hợp các agent chạy đồng thời | Nhiều agent cần giao tiếp, chia sẻ state với nhau |
| **`performance-monitor`** | Hạ tầng observability, theo dõi metric | Cần theo dõi hiệu năng hệ thống, phát hiện anomaly |
| **`task-distributor`** | Phân phối task, quản lý queue, cân bằng tải | Cần phân chia công việc cho nhiều agent/worker |
| **`workflow-orchestrator`** | Thiết kế workflow nghiệp vụ | Cần tự động hóa quy trình nhiều bước với error handling |

### 🛡 Quality & Security (16 agents)

Các agent cho kiểm thử, bảo mật, và đảm bảo chất lượng.

| Agent | Mô tả | Sử dụng khi nào |
|---|---|---|
| **`accessibility-tester`** | Tuân thủ WCAG, kiểm thử accessibility | Cần kiểm tra/sửa lỗi accessibility, đạt chuẩn WCAG |
| **`ad-security-reviewer`** | Đánh giá bảo mật Active Directory | Cần audit AD, kiểm tra privilege escalation |
| **`ai-writing-auditor`** | Kiểm tra nội dung do AI tạo | Cần phát hiện/viết lại văn bản AI-generated |
| **`architect-reviewer`** | Review thiết kế hệ thống | Cần đánh giá architectural decision, review design |
| **`chaos-engineer`** | Thí nghiệm failure có kiểm soát | Cần test resilience, inject failure, hoặc game day |
| **`code-reviewer`** | Review code toàn diện — bảo mật, chất lượng | Cần review PR, kiểm tra code quality hoặc security |
| **`compliance-auditor`** | Tuân thủ quy định, kiểm soát audit | Cần đạt chuẩn SOC2, GDPR, hoặc ISO |
| **`debugger`** | Chẩn đoán bug, phân tích nguyên nhân gốc | Gặp bug khó, cần root cause analysis |
| **`error-detective`** | Tương quan lỗi, phân tích chuỗi failure | Lỗi liên quan nhiều service, cần trace failure chain |
| **`penetration-tester`** | Kiểm thử xâm nhập bảo mật | Cần pentest, tìm vulnerability trước khi release |
| **`performance-engineer`** | Xác định nút cổ chai hiệu năng | Ứng dụng chậm, cần profiling hoặc load test |
| **`powershell-security-hardening`** | Bảo mật PowerShell, hardening | Cần hardening PS remoting, constrained mode |
| **`qa-expert`** | Chiến lược QA, kế hoạch test, coverage | Cần lập chiến lược test, đánh giá coverage |
| **`security-auditor`** | Kiểm toán bảo mật, đánh giá tuân thủ | Cần audit toàn diện, đánh giá rủi ro bảo mật |
| **`test-automator`** | Framework test tự động, CI/CD | Cần xây test framework, tích hợp testing vào CI |
| **`ui-ux-tester`** | Kiểm thử chức năng UI/UX | Cần test giao diện, user flow, tìm defect |

### 🔬 Research & Analysis (8 agents)

Các agent cho nghiên cứu, phân tích, và thu thập insight.

| Agent | Mô tả | Sử dụng khi nào |
|---|---|---|
| **`competitive-analyst`** | Phân tích đối thủ, benchmarking | Cần so sánh với đối thủ, đánh giá vị thế thị trường |
| **`data-researcher`** | Khám phá và xác thực dữ liệu đa nguồn | Cần thu thập dữ liệu từ nhiều nguồn, xác thực độ tin cậy |
| **`market-researcher`** | Phân tích thị trường, hành vi người tiêu dùng | Cần hiểu thị trường, khách hàng mục tiêu, hoặc sizing |
| **`project-idea-validator`** | Stress-test ý tưởng, phân tích đối thủ | Có ý tưởng mới, cần kiểm tra tính khả thi trước khi đầu tư |
| **`research-analyst`** | Tổng hợp nghiên cứu đa nguồn | Cần tổng hợp thông tin từ nhiều nguồn thành báo cáo |
| **`scientific-literature-researcher`** | Tài liệu khoa học, dữ liệu cấu trúc | Cần tìm/trích xuất dữ liệu từ bài báo khoa học |
| **`search-specialist`** | Chiến lược tìm kiếm nâng cao, tối ưu query | Cần tìm thông tin chính xác nhanh từ nhiều nguồn |
| **`trend-analyst`** | Xu hướng mới nổi, dự đoán biến đổi ngành | Cần dự báo xu hướng, lập kịch bản tương lai |

### 🌐 Specialized Domains (13 agents)

Các agent cho các lĩnh vực chuyên biệt.

| Agent | Mô tả | Sử dụng khi nào |
|---|---|---|
| **`api-documenter`** | Tài liệu API, đặc tả OpenAPI | Cần viết/cải thiện API docs, Swagger spec |
| **`blockchain-developer`** | Smart contract, DApp, blockchain | Xây smart contract, DApp, hoặc DeFi protocol |
| **`embedded-systems`** | Firmware, RTOS, vi điều khiển | Lập trình embedded, firmware, hoặc RTOS |
| **`fintech-engineer`** | Hệ thống thanh toán, tuân thủ tài chính | Xây payment system, PSD2, hoặc banking API |
| **`game-developer`** | Game system, đồ họa, multiplayer | Xây game, rendering, physics, hoặc netcode |
| **`healthcare-admin`** | Quản trị y tế, tuân thủ HIPAA | Dự án healthcare, EMR, hoặc cần HIPAA compliance |
| **`iot-engineer`** | Quản lý thiết bị IoT, edge computing | Xây IoT platform, device management, edge |
| **`m365-admin`** | Tự động hóa quản trị Microsoft 365 | Cần tự động Exchange, SharePoint, Teams |
| **`mobile-app-developer`** | Ứng dụng iOS/Android | Xây mobile app native hoặc cross-platform |
| **`payment-integration`** | Tích hợp cổng thanh toán, PCI | Tích hợp Stripe, PayPal, hoặc cần PCI compliance |
| **`quant-analyst`** | Giao dịch định lượng, mô hình tài chính | Xây trading strategy, risk model, hoặc backtesting |
| **`risk-manager`** | Nhận diện/giảm thiểu rủi ro | Cần đánh giá rủi ro, thiết kế control framework |
| **`seo-specialist`** | Kiểm toán SEO, chiến lược từ khóa | Cần tối ưu SEO, technical audit, hoặc keyword research |

---

## 🚀 Cách sử dụng

### Bật agent/skill trong config

```yaml
# .dotagen/config.yaml
agents:
  dotagent-game-studios-lead-programmer:
    targets: all          # Bật cho tất cả nền tảng
  dotagent-voltagent-code-reviewer:
    targets:
      - claude-code       # Chỉ Claude Code

skills:
  dotagent-game-studios-gate-check:
    targets: all
  dotagent-superpowers-test-driven-development:
    targets:
      - claude-code
      - cursor
```

### Đồng bộ

```bash
dotagen sync              # Đồng bộ tất cả
dotagen sync claude-code  # Đồng bộ chỉ Claude Code
```

### Quản lý qua CLI

```bash
dotagen skill list                 # Liệt kê skills
dotagen skill create my-workflow   # Tạo skill mới
dotagen status                     # Kiểm tra trạng thái
```

### Quản lý qua Web Dashboard

```bash
dotagen serve   # Mở dashboard tại http://localhost:7890
```

---

## 📝 Nguồn gốc

- **49 agents + 73 skills** Game Studios lấy từ [donchitos/claude-code-game-studios](https://github.com/donchitos/claude-code-game-studios) (MIT)
- **154 agents** `dotagent-voltagent-*` lấy từ [VoltAgent/awesome-claude-code-subagents](https://github.com/VoltAgent/awesome-claude-code-subagents)
- **26 skills** `dotagent-mattpocock-*` lấy từ [mattpocock/skills](https://github.com/mattpocock/skills)
- **14 skills** `dotagent-superpowers-*` lấy từ [obra/superpowers](https://github.com/obra/superpowers)
- Các skill vendor còn lại lấy từ registry [awesome-agent-skills](https://github.com/enolalabs/awesome-agent-skills); [Hallmark](https://github.com/Nutlope/hallmark) là snapshot MIT riêng
