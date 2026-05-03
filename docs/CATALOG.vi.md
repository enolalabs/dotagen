# 📚 Danh mục Agents & Skills mặc định

> Tài liệu giới thiệu toàn bộ **144 agents** và **16 skills** được tích hợp sẵn trong dotagen.
> Tất cả đều bị tắt theo mặc định — bạn chọn bật agent/skill nào và cho nền tảng nào trong `config.yaml`.

---

## 🛠 Skills (Slash Commands)

Skills là các **quy trình làm việc có cấu trúc** (slash command) mà agent sẽ tuân theo khi được kích hoạt. Khác với agent (định nghĩa *agent là ai*), skill định nghĩa *agent làm gì theo quy trình nào*.

Mỗi skill được lưu dưới dạng thư mục `ds-<tên>/SKILL.md`, có thể kèm thư mục `references/` chứa các file tham chiếu bổ sung.

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

## 🤖 Agents

Agents là các **chuyên gia ảo** với vai trò và chuyên môn cụ thể. Mỗi agent chứa system prompt hướng dẫn AI hành xử theo đúng vai trò khi được kích hoạt.

Mỗi agent được lưu dưới dạng file `da-<tên>.md` với frontmatter YAML (description, category) và nội dung Markdown.

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
  da-backend-developer:
    targets: all          # Bật cho tất cả nền tảng
  da-code-reviewer:
    targets:
      - claude-code       # Chỉ Claude Code

skills:
  ds-diagnose:
    targets: all
  ds-tdd:
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

- **144 agents** được lấy từ dự án [VoltAgent](https://github.com/VoltAgent/voltagent)
- **16 skills** được lấy từ [mattpocock/skills](https://github.com/mattpocock/skills)
