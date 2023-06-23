## 快递代取

默认端口：9090

### 订单相关接口
**订单列表：** `/express/list` Get

**订单详情：** `/express/info` Get

**新增订单：** `/express/create` Post

**接单：** `/express/order` Put

**完成：** `/express/finish` Put

### 用户相关接口
**用户登录：** `/user/login` Post
>用户存在时只进行登录，用户不存在时注册并登录

**用户详情：** `/user/info` Get

**修改用户信息：** `/user/info` Put

### 公用接口
**上传文件：** `/upload` Post

