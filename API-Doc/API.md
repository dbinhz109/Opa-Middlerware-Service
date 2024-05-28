# ACTIVE API

## 📦 API

### APP

| Name | Description | Flow |
| --- | --- | --- |
| get-one/:id | Tìm thông tin của 1 app theo id trong database và trả về giao diện | Giao diện -> lsmanagement -> database -> giao diện | 
| find-all | Tìm thông tin tất cả các app trong database và trả về giao diện | Giao diện -> lsmanagement -> database -> giao diện |

### APP CATEGORY

| Name | Description | Flow |
| --- | --- | --- |
| get-one/:id | Tìm thông tin của 1 app category theo id trong database và trả về giao diện | Giao diện -> lsmanagement -> database -> giao diện | 
| find-all | Tìm thông tin tất cả các app category trong database và trả về giao diện | Giao diện -> lsmanagement -> database -> giao diện |

### WEB

| Name | Description | Flow |
| --- | --- | --- |
| get-one/:id | Tìm thông tin của 1 web theo id trong database và trả về giao diện | Giao diện -> lsmanagement -> database -> giao diện | 
| find-all | Tìm thông tin tất cả các web trong database và trả về giao diện | Giao diện -> lsmanagement -> database -> giao diện |

### WEB CATEGORY

| Name | Description | Flow |
| --- | --- | --- |
| get-one/:id | Tìm thông tin của 1 web category theo id trong database và trả về giao diện | Giao diện -> lsmanagement -> database -> giao diện | 
| find-all | Tìm thông tin tất cả các web category trong database và trả về giao diện | Giao diện -> lsmanagement -> database -> giao diện |

### PROTOCOL

| Name | Description | Flow |
| --- | --- | --- |
| get-one/:id | Tìm thông tin của 1 protocol theo id trong database và trả về giao diện | Giao diện -> lsmanagement -> database -> giao diện | 
| find-all | Tìm thông tin tất cả các protocol trong database và trả về giao diện | Giao diện -> lsmanagement -> database -> giao diện |

### PROTOCOL CATGEGORY

| Name | Description | Flow |
| --- | --- | --- |
| get-one/:id | Tìm thông tin của 1 protocol category theo id trong database và trả về giao diện | Giao diện -> lsmanagement -> database -> giao diện | 
| find-all | Tìm thông tin tất cả các protocol category trong database và trả về giao diện | Giao diện -> lsmanagement -> database -> giao diện |

### USER

| Name | Description | Flow |
| --- | --- | --- |
| create | Tạo mới user lưu vào trong database và gửi yêu cầu cho lscontroller | Giao diện -> lsmanagement -> database -> lscontroller | 
| update/:id | Cập nhập thông tin user theo id lưu vào trong database và gửi yêu cầu cho lscontroller | Giao diện -> lsmanagement -> database -> lscontroller |
| delete/:id | Xóa user khỏi database theo id và gửi yêu cầu cho lscontroller | Giao diện -> lsmanagement -> database -> lscontroller |
| get-one/:id | Tìm thông tin của 1 user theo id trong database và trả về giao diện | Giao diện -> lsmanagement -> database -> giao diện | 
| get-all | Tìm thông tin tất cả các user trong database và trả về giao diện | Giao diện -> lsmanagement -> database -> giao diện |

### End Device

| Name | Description | Flow |
| --- | --- | --- |
| delete/:id | Xóa end device khỏi database theo id và gửi yêu cầu cho lscontroller | Giao diện -> lsmanagement -> database -> lscontroller |
| get-one/:id | Tìm thông tin của 1 user theo id trong database và trả về giao diện | Giao diện -> lsmanagement -> database -> giao diện| 
| get-all/:userId | Tìm thông tin tất cả các user theo userid trong database và trả về giao diện | Giao diện -> lsmanagement -> database -> giao diện |
| find-all-by-adminId | Tìm thông tin tất cả các user theo adminid trong database và trả về giao diện | Giao diện -> lsmanagement -> database -> giao diện |

### Edges

| Name | Description | Flow |
| --- | --- | --- |
| delete/:id | Xóa edge khỏi database theo id và gửi yêu cầu cho lscontroller | Giao diện -> lsmanagement -> database -> lscontroller |
| get-all | Tìm thông tin tất cả các edge trong database và gửi yêu cầu cho lscontroller | Giao diện -> lsmanagement -> database -> giao diện |
| get-all-by-accountType | Tìm thông tin tất cả các edge theo account type trong database và trả về giao diện | Giao diện -> lsmanagement -> database -> giao diện |
| get-edge-by-serial/:serial | Tìm thông tin edge theo serial trong database và trả về giao diện | Giao diện -> lsmanagement -> database -> giao diện |
**| find-all-edge-apply-scenario-private/:scenarioId | * | Giao diện -> lsmanagement -> database -> giao diện |
**| deregister-edge | * | Giao diện -> lsmanagement -> database -> lscontroller |
**| register-edge | * | Giao diện -> lsmanagement -> database -> lscontroller |

### Scenario

| Name | Description | Flow |
| --- | --- | --- |
| delete/:id | Xóa scenario khỏi database theo id và gửi yêu cầu cho lscontroller | Giao diện -> lsmanagement -> database -> lscontroller |
| create | Tạo mới scenario lưu vào trong database và gửi yêu cầu cho lscontroller | Giao diện -> lsmanagement -> database -> lscontroller |
| update/:id | Cập nhập thông tin scenario lưu vào trong database và gửi yêu cầu cho lscontroller | Giao diện -> lsmanagement -> database -> lscontroller |
| get-one/:id | Tìm thông tin của 1 scenario theo id trong database và trả về giao diện | Giao diện -> lsmanagement -> database -> giao diện |
| get-all | Tìm thông tin tất cả các edge trong database và trả về giao diện | Giao diện -> lsmanagement -> database -> giao diện |
**| find-scenario-apply-global | * | Giao diện -> lsmanagement -> database -> giao diện |
**| find-scenario-apply-private/:edgeID | * | Giao diện -> lsmanagement -> database -> giao diện |
**| status-apply-scenario-for-edge/:edgeID | * | Giao diện -> lsmanagement -> database -> giao diện |
**| scenario-config-apply | * | Giao diện -> lsmanagement -> database -> giao diện |
**| dismiss-scenario-config-apply | * | Giao diện -> lsmanagement -> database -> giao diện |
**| scenario-config-apply-all | * | Giao diện -> lsmanagement -> database -> giao diện |

### Policy

| Name | Description | Flow |
| --- | --- | --- |
| delete/:id | Xóa policy khỏi database theo id và gửi yêu cầu cho lscontroller | Giao diện -> lsmanagement -> database -> lscontroller |
| create | Tạo mới policy lưu vào trong database và gửi yêu cầu cho lscontroller | Giao diện -> lsmanagement -> database -> lscontroller |
| update/:id | Cập nhập thông tin policy lưu vào trong database và gửi yêu cầu cho lscontroller | Giao diện -> lsmanagement -> database -> lscontroller |
| get-one/:id | Tìm thông tin của 1 policy theo id trong database và trả về giao diện | Giao diện -> lsmanagement -> database -> giao diện |
| get-policy-in-scenario/:scenarioId | Tìm thông tin của policy trong scenario trong database và trả về giao diện | Giao diện -> lsmanagement -> database -> giao diện |
| delete-by-mac/:mac | Xóa policy khỏi database theo mac và gửi yêu cầu cho lscontroller | Giao diện -> lsmanagement -> database -> lscontroller |
| delete-by-userid/:id | Xóa policy khỏi database theo userid và gửi yêu cầu cho lscontroller | Giao diện -> lsmanagement -> database -> lscontroller |

### Network

| Name | Description | Flow |
| --- | --- | --- |
| show-network-interface/:edgeID | Hiển thị network interface theo edge ID và gửi yêu cầu cho lscontroller rồi trả về giao diện | Giao diện -> lsmanagement -> lscontroller -> giao diện |
| get-uuid/:edgeID | Lấy uuid từ dưới edge bằng edge ID và gửi yêu cầu cho lscontroller rồi trả về giao diện | Giao diện -> lsmanagement -> lscontroller -> giao diện |
| show-interface-wan/:edgeID | Hiển thị interface wan theo edge ID và gửi yêu cầu cho lscontroller rồi trả về giao diện | Giao diện -> lsmanagement -> lscontroller -> giao diện |
| show-interface-lan/:edgeID | Hiển thị interface lan theo edge ID và gửi yêu cầu cho lscontroller rồi trả về giao diện | Giao diện -> lsmanagement -> lscontroller -> giao diện |
| config-network-lan | Thay đổi cấu hình interface lan và gửi yêu cầu cho lscontroller | Giao diện -> lsmanagement -> lscontroller |
| config-network-wan | Thay đổi cấu hình interface wan và gửi yêu cầu cho lscontroller | Giao diện -> lsmanagement -> lscontroller |

### Authentication

| Name | Description | Flow |
| --- | --- | --- |
| config-auth | Thay đổi cấu hình xác thực thiết bị và gửi yêu cầu cho lscontroller | Giao diện -> lsmanagement -> database -> lscontroller |
| config-deauth | Thay đổi cấu hình xác thực thiết bị và gửi yêu cầu cho lscontroller | Giao diện -> lsmanagement -> database -> lscontroller |

### Basic wifi

| Name | Description | Flow |
| --- | --- | --- |
| show-network-wireless/:edgeID | Hiển thị network wireless theo edge ID và gửi yêu cầu cho lscontroller rồi trả về giao diện | Giao diện -> lsmanagement -> lscontroller -> giao diện |
| show-network-guest/:edgeID | Hiển thị interface guest theo edge ID và gửi yêu cầu cho lscontroller rồi trả về giao diện | Giao diện -> lsmanagement -> lscontroller -> giao diện |
| config-network-wireless | Thay đổi cấu hình network wireless và gửi yêu cầu cho lscontroller | Giao diện -> lsmanagement -> lscontroller |
| config-network-guest | Thay đổi cấu hình network guest và gửi yêu cầu cho lscontroller | Giao diện -> lsmanagement -> lscontroller |
| config-mesh-controller | Thay đổi cấu hình mesh controller và gửi yêu cầu cho lscontroller | Giao diện -> lsmanagement -> lscontroller |
| config-mesh-agent | Thay đổi cấu hình mesh agent và gửi yêu cầu cho lscontroller | Giao diện -> lsmanagement -> lscontroller |

### System

| Name | Description | Flow |
| --- | --- | --- |
| show-system-info/:edgeID | Hiển thị system info theo edge ID và gửi yêu cầu cho lscontroller rồi trả về giao diện | Giao diện -> lsmanagement -> lscontroller -> giao diện |
| show-system-board/:edgeID | Hiển thị system board theo edge ID và gửi yêu cầu cho lscontroller rồi trả về giao diện | Giao diện -> lsmanagement -> lscontroller -> giao diện |
| config-system | Thay đổi cấu hình thiết bị và gửi yêu cầu cho lscontroller | Giao diện -> lsmanagement -> lscontroller |
| config-system-update | Cập nhập thiết bị và gửi yêu cầu cho lscontroller | Giao diện -> lsmanagement -> lscontroller |
| config-system-reboot | Khởi động lại thiết bị và gửi yêu cầu cho lscontroller | Giao diện -> lsmanagement -> lscontroller |
| config-system-backup | Sao lưu dữ liệu trên thiết bị và gửi yêu cầu cho lscontroller | Giao diện -> lsmanagement -> lscontroller |
| config-system-factory | Đặt lại cấu hình nhà sản xuất và gửi yêu cầu cho lscontroller | Giao diện -> lsmanagement -> lscontroller |

### AAR






