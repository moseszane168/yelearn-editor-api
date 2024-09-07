/*
 Navicat Premium Data Transfer

 Source Server         : dev
 Source Server Type    : MySQL
 Source Server Version : 50732
 Source Host           : 172.16.1.66:3306
 Source Schema         : crf_mold

 Target Server Type    : MySQL
 Target Server Version : 50732
 File Encoding         : 65001

 Date: 28/06/2022 16:39:29
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for dict_group
-- ----------------------------
DROP TABLE IF EXISTS `dict_group`;
CREATE TABLE `dict_group`  (
  `id` bigint(64) NOT NULL AUTO_INCREMENT,
  `code` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '组编码',
  `name` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '组名称',
  `is_deleted` varchar(1) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT 'N' COMMENT '是否删除',
  `gmt_created` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '创建时间',
  `gmt_updated` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) ON UPDATE CURRENT_TIMESTAMP(0) COMMENT '修改时间',
  `created_by` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '创建人',
  `updated_by` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '修改人',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 64 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '字典组' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for dict_property
-- ----------------------------
DROP TABLE IF EXISTS `dict_property`;
CREATE TABLE `dict_property`  (
  `id` bigint(64) NOT NULL AUTO_INCREMENT,
  `key` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT 'key',
  `value_cn` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT 'value_cn',
  `value_en` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT 'value_en',
  `order` int(255) UNSIGNED NOT NULL COMMENT '排序',
  `group_code` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '组编码',
  `is_deleted` varchar(1) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT 'N' COMMENT '是否删除',
  `gmt_created` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '创建时间',
  `gmt_updated` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) ON UPDATE CURRENT_TIMESTAMP(0) COMMENT '修改时间',
  `created_by` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '创建人',
  `updated_by` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '修改人',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 139 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '字典' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for login_info
-- ----------------------------
DROP TABLE IF EXISTS `login_info`;
CREATE TABLE `login_info`  (
  `id` bigint(64) NOT NULL AUTO_INCREMENT,
  `login_name` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '工号',
  `last_login_date` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0) ON UPDATE CURRENT_TIMESTAMP(0) COMMENT '上次登录时间',
  `fault_count` int(11) NOT NULL DEFAULT 0 COMMENT '登录错误次数',
  `is_deleted` varchar(1) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT 'N' COMMENT '是否删除',
  `gmt_created` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '创建时间',
  `gmt_updated` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) ON UPDATE CURRENT_TIMESTAMP(0) COMMENT '修改时间',
  `created_by` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '创建人',
  `updated_by` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '修改人',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 28 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '登录信息' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for message
-- ----------------------------
DROP TABLE IF EXISTS `message`;
CREATE TABLE `message`  (
  `id` bigint(64) NOT NULL AUTO_INCREMENT,
  `status` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '状态',
  `content` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '消息内容',
  `job_id` bigint(64) NULL DEFAULT NULL COMMENT '任务ID',
  `remodel_id` bigint(64) NULL DEFAULT NULL COMMENT '改造ID',
  `type` varchar(60) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '消息类型：maintenance、remodel',
  `is_deleted` varchar(1) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT 'N' COMMENT '是否删除',
  `gmt_created` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '创建时间',
  `gmt_updated` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) ON UPDATE CURRENT_TIMESTAMP(0) COMMENT '修改时间',
  `created_by` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '创建人',
  `updated_by` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '修改人',
  `operator` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '邮件接收人',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1367437 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '消息' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for mold_bom
-- ----------------------------
DROP TABLE IF EXISTS `mold_bom`;
CREATE TABLE `mold_bom`  (
  `id` bigint(64) NOT NULL AUTO_INCREMENT,
  `standard_component` varchar(1) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT 'Y' COMMENT '标准件',
  `mold_code` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '模具编号',
  `part_name` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '零件代号',
  `part_code` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '零件名称',
  `count` int(11) NOT NULL COMMENT '数量',
  `material` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '材料',
  `flavor` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '规格',
  `standard` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '标准',
  `heat_treating` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '热处理',
  `assembly_relation` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '装配关系',
  `remark` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '备注',
  `is_deleted` varchar(1) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT 'N' COMMENT '是否删除',
  `created_by` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '创建人',
  `updated_by` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '修改人',
  `gmt_created` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '创建时间',
  `gmt_updated` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) ON UPDATE CURRENT_TIMESTAMP(0) COMMENT '修改时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 25 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '模具bom' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for mold_custom_info
-- ----------------------------
DROP TABLE IF EXISTS `mold_custom_info`;
CREATE TABLE `mold_custom_info`  (
  `id` bigint(64) NOT NULL AUTO_INCREMENT,
  `mold_code` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '模具ID',
  `key` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '字段名',
  `value` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '字段值',
  `is_deleted` varchar(1) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT 'N' COMMENT '是否删除',
  `created_by` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '创建人',
  `updated_by` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '修改人',
  `gmt_created` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '创建时间',
  `gmt_updated` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) ON UPDATE CURRENT_TIMESTAMP(0) COMMENT '修改时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 49 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '台账自定义字段' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for mold_doc
-- ----------------------------
DROP TABLE IF EXISTS `mold_doc`;
CREATE TABLE `mold_doc`  (
  `id` bigint(64) NOT NULL AUTO_INCREMENT,
  `file_key` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '文件Minio key',
  `mold_code` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '模具编号',
  `name` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '文档名称',
  `version` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '版本',
  `content` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '变更内容',
  `remark` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '备注',
  `is_deleted` varchar(1) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT 'N' COMMENT '是否删除',
  `created_by` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '创建人',
  `updated_by` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '修改人',
  `gmt_created` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '创建时间',
  `gmt_updated` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) ON UPDATE CURRENT_TIMESTAMP(0) COMMENT '修改时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 5 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '模具文档' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for mold_info
-- ----------------------------
DROP TABLE IF EXISTS `mold_info`;
CREATE TABLE `mold_info`  (
  `id` bigint(64) NOT NULL AUTO_INCREMENT,
  `code` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '模具编号',
  `name` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '模具名称',
  `client_name` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '客户名称',
  `weight` double NOT NULL COMMENT '模具重量',
  `project_name` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '项目名称',
  `size_long` double NULL DEFAULT NULL COMMENT '模具长',
  `size_width` double NULL DEFAULT NULL COMMENT '模具宽',
  `size_heigh` double NULL DEFAULT NULL COMMENT '模具高',
  `platform` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '平台',
  `process` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '工序',
  `property_number` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '资产编号',
  `type` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '模具类型',
  `status` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '模具状态',
  `make_date` datetime(0) NULL DEFAULT NULL COMMENT '制造日期',
  `first_use_date` datetime(0) NULL DEFAULT NULL COMMENT '第一次使用日期',
  `year_limit` double(6, 2) NULL DEFAULT NULL COMMENT '使用年限',
  `rfid` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT 'RFID',
  `line_level` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '使用线别',
  `provide` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '模具供应商',
  `category` varchar(2) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '模具类别：A、B',
  `flush_count` bigint(64) NULL DEFAULT NULL COMMENT '模具冲次',
  `is_deleted` varchar(1) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT 'N' COMMENT '是否删除',
  `created_by` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '创建人',
  `updated_by` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '修改人',
  `gmt_created` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '创建时间',
  `gmt_updated` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) ON UPDATE CURRENT_TIMESTAMP(0) COMMENT '修改时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1235 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '模具台账' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for mold_knowledge
-- ----------------------------
DROP TABLE IF EXISTS `mold_knowledge`;
CREATE TABLE `mold_knowledge`  (
  `id` bigint(64) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '知识名称',
  `group` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '知识分类',
  `content` mediumtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '知识内容',
  `style` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '书皮样式',
  `is_deleted` varchar(1) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT 'N' COMMENT '是否删除',
  `gmt_created` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '创建时间',
  `gmt_updated` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) ON UPDATE CURRENT_TIMESTAMP(0) COMMENT '修改时间',
  `created_by` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '创建人',
  `updated_by` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '修改人',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 42 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '模具知识库' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for mold_maintenance_plan
-- ----------------------------
DROP TABLE IF EXISTS `mold_maintenance_plan`;
CREATE TABLE `mold_maintenance_plan`  (
  `id` bigint(64) NOT NULL AUTO_INCREMENT,
  `code` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '计划编号',
  `name` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '模具保养计划名称',
  `mold_type` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '模具类型',
  `task_start` int(11) NOT NULL COMMENT '任务生成区间最小值',
  `task_end` int(11) NOT NULL COMMENT '任务生成区间最大值，大于',
  `task_standard` int(11) NOT NULL COMMENT '任务生成标准值，只做显示用',
  `operate` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '任务生成区间最小值操作符：小于等于lte、小于lt、等于ge',
  `status` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '任务状态：running、pause',
  `is_deleted` varchar(1) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT 'N' COMMENT '是否删除',
  `created_by` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '创建人',
  `updated_by` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '修改人',
  `gmt_created` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '创建时间',
  `gmt_updated` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) ON UPDATE CURRENT_TIMESTAMP(0) COMMENT '修改时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 10 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '模具保养计划' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for mold_maintenance_standard
-- ----------------------------
DROP TABLE IF EXISTS `mold_maintenance_standard`;
CREATE TABLE `mold_maintenance_standard`  (
  `id` bigint(64) NOT NULL AUTO_INCREMENT,
  `name` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '保养名称',
  `code` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '标准编号',
  `version` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '版本',
  `type` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '模具类型',
  `level` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '保养级别',
  `standard_type` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '标准类型',
  `is_deleted` varchar(1) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT 'N' COMMENT '是否删除',
  `created_by` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '创建人',
  `updated_by` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '修改人',
  `gmt_created` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '创建时间',
  `gmt_updated` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) ON UPDATE CURRENT_TIMESTAMP(0) COMMENT '修改时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 15 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '模具保养标准' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for mold_maintenance_standard_content
-- ----------------------------
DROP TABLE IF EXISTS `mold_maintenance_standard_content`;
CREATE TABLE `mold_maintenance_standard_content`  (
  `id` bigint(64) NOT NULL AUTO_INCREMENT,
  `mold_maintenance_standard_id` bigint(64) NOT NULL COMMENT '模具保养标准表ID',
  `cycle` int(11) NOT NULL COMMENT '周期',
  `item` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '维护项',
  `judge_method` varchar(300) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '维护判断方法',
  `maintenance_method` varchar(300) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '维护方法',
  `eligibility_criteria` varchar(300) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '合格标准',
  `remark` varchar(300) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '备注',
  `order` int(11) NULL DEFAULT NULL COMMENT '顺序',
  `is_deleted` varchar(1) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT 'N' COMMENT '是否删除',
  `created_by` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '创建人',
  `updated_by` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '修改人',
  `gmt_created` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '创建时间',
  `gmt_updated` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) ON UPDATE CURRENT_TIMESTAMP(0) COMMENT '修改时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 38 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '模具保养标准内容' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for mold_maintenance_standard_rel
-- ----------------------------
DROP TABLE IF EXISTS `mold_maintenance_standard_rel`;
CREATE TABLE `mold_maintenance_standard_rel`  (
  `id` bigint(64) NOT NULL AUTO_INCREMENT,
  `mold_maintenance_standard_id` bigint(64) NOT NULL COMMENT '模具保养标准',
  `mold_id` bigint(64) NOT NULL COMMENT '模具ID',
  `is_deleted` varchar(1) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT 'N' COMMENT '是否删除',
  `created_by` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '创建人',
  `updated_by` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '修改人',
  `gmt_created` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '创建时间',
  `gmt_updated` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) ON UPDATE CURRENT_TIMESTAMP(0) COMMENT '修改时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 44 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '模具保养标准关联表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for mold_maintenance_task
-- ----------------------------
DROP TABLE IF EXISTS `mold_maintenance_task`;
CREATE TABLE `mold_maintenance_task`  (
  `id` bigint(64) UNSIGNED NOT NULL AUTO_INCREMENT,
  `code` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '任务编号',
  `mold_id` bigint(64) NOT NULL COMMENT '模具编号',
  `mold_maintenance_plain_id` bigint(64) NULL DEFAULT NULL COMMENT '模具标准计划ID',
  `mold_maintenance_standard_id` bigint(64) NULL DEFAULT NULL COMMENT '模具标准ID',
  `status` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '状态',
  `standard_name` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '标准名称',
  `maintenance_level` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '保养级别',
  `reason` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '挂起原因',
  `operator` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '保养人',
  `time` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '保养时间',
  `hang_up_count` int(11) NULL DEFAULT 0 COMMENT '挂起次数',
  `task_gen_interval` int(11) NULL DEFAULT NULL COMMENT '任务生成区间',
  `timeout_count` bigint(64) NULL DEFAULT NULL COMMENT '超时冲次',
  `is_deleted` varchar(1) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT 'N' COMMENT '是否删除',
  `created_by` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '创建人',
  `updated_by` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '修改人',
  `gmt_created` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '创建时间',
  `gmt_updated` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) ON UPDATE CURRENT_TIMESTAMP(0) COMMENT '修改时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 52119 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '模具保养任务' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for mold_maintenance_task_content
-- ----------------------------
DROP TABLE IF EXISTS `mold_maintenance_task_content`;
CREATE TABLE `mold_maintenance_task_content`  (
  `id` bigint(64) NOT NULL AUTO_INCREMENT,
  `mold_maintenance_task_id` bigint(64) NOT NULL COMMENT '模具保养任务ID',
  `cycle` int(11) NOT NULL COMMENT '周期',
  `item` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '维护项',
  `judge_method` varchar(300) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '维护判断方法',
  `maintenance_method` varchar(300) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '维护方法',
  `eligibility_criteria` varchar(300) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '合格标准',
  `remark` varchar(300) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '备注',
  `order` int(11) NULL DEFAULT NULL COMMENT '顺序',
  `status` varchar(1) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '状态',
  `is_deleted` varchar(1) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT 'N' COMMENT '是否删除',
  `created_by` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '创建人',
  `updated_by` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '修改人',
  `gmt_created` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '创建时间',
  `gmt_updated` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) ON UPDATE CURRENT_TIMESTAMP(0) COMMENT '修改时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 34 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '模具保养标准内容' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for mold_maintenance_task_lifecycle
-- ----------------------------
DROP TABLE IF EXISTS `mold_maintenance_task_lifecycle`;
CREATE TABLE `mold_maintenance_task_lifecycle`  (
  `id` bigint(64) UNSIGNED NOT NULL AUTO_INCREMENT,
  `mold_maintenance_task_id` bigint(64) NOT NULL COMMENT '模具标准任务ID',
  `title` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '标题,字典key',
  `operator` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '操作人',
  `time` datetime(0) NOT NULL COMMENT '时间',
  `content` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '内容',
  `is_deleted` varchar(1) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT 'N' COMMENT '是否删除',
  `created_by` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '创建人',
  `updated_by` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '修改人',
  `gmt_created` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '创建时间',
  `gmt_updated` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) ON UPDATE CURRENT_TIMESTAMP(0) COMMENT '修改时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 48340 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '模具保养任务生命周期（任务时间线）' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for mold_params
-- ----------------------------
DROP TABLE IF EXISTS `mold_params`;
CREATE TABLE `mold_params`  (
  `id` bigint(64) NOT NULL AUTO_INCREMENT,
  `mold_code` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '模具编号',
  `height` double NOT NULL COMMENT '模具高度',
  `defensive_device` varchar(1) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '防错装置',
  `desc` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '描述',
  `is_deleted` varchar(1) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT 'N' COMMENT '是否删除',
  `created_by` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '创建人',
  `updated_by` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '修改人',
  `gmt_created` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '创建时间',
  `gmt_updated` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) ON UPDATE CURRENT_TIMESTAMP(0) COMMENT '修改时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 5 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '模具成型参数' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for mold_params_custom
-- ----------------------------
DROP TABLE IF EXISTS `mold_params_custom`;
CREATE TABLE `mold_params_custom`  (
  `id` bigint(64) NOT NULL AUTO_INCREMENT,
  `mold_code` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '模具ID',
  `key` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '字段名',
  `value` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '字段值',
  `is_deleted` varchar(1) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT 'N' COMMENT '是否删除',
  `created_by` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '创建人',
  `updated_by` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '修改人',
  `gmt_created` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '创建时间',
  `gmt_updated` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) ON UPDATE CURRENT_TIMESTAMP(0) COMMENT '修改时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 53 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '模具成型参数自定义字段' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for mold_part_rel
-- ----------------------------
DROP TABLE IF EXISTS `mold_part_rel`;
CREATE TABLE `mold_part_rel`  (
  `id` bigint(64) NOT NULL AUTO_INCREMENT,
  `mold_code` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '模具code',
  `part_code` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '零件号',
  `is_deleted` varchar(1) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT 'N' COMMENT '是否删除',
  `created_by` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '创建人',
  `updated_by` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '修改人',
  `gmt_created` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '创建时间',
  `gmt_updated` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) ON UPDATE CURRENT_TIMESTAMP(0) COMMENT '修改时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 128 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '模具零件' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for mold_product_resume
-- ----------------------------
DROP TABLE IF EXISTS `mold_product_resume`;
CREATE TABLE `mold_product_resume`  (
  `id` bigint(64) NOT NULL AUTO_INCREMENT,
  `mold_code` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '模具编码，冗余存储',
  `mold_name` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '模具名称，冗余存储',
  `mold_type` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '模具类型，冗余存储',
  `mold_id` bigint(64) NOT NULL COMMENT '模具ID',
  `order_code` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '工单号',
  `part_code` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '零件号',
  `lineLevel` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '线别',
  `complete_time` datetime(0) NOT NULL COMMENT '生产完成时间',
  `count` bigint(64) NOT NULL COMMENT '工单数量',
  `is_deleted` varchar(1) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT 'N' COMMENT '是否删除',
  `is_checked` varchar(1) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT 'Y' COMMENT '是否已经检查过',
  `created_by` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '创建人',
  `updated_by` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '修改人',
  `gmt_created` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '创建时间',
  `gmt_updated` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) ON UPDATE CURRENT_TIMESTAMP(0) COMMENT '修改时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1494 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '模具生产履历' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for mold_remodel
-- ----------------------------
DROP TABLE IF EXISTS `mold_remodel`;
CREATE TABLE `mold_remodel`  (
  `id` bigint(64) NOT NULL AUTO_INCREMENT,
  `code` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '改造编号',
  `mold_code` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '模具编码',
  `remodel_start_time` datetime(0) NOT NULL COMMENT '改造开始时间',
  `remodel_end_time` datetime(0) NOT NULL COMMENT '改造结束时间',
  `finish_time` datetime(0) NULL DEFAULT NULL COMMENT '改造完成时间',
  `type` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '改造类别',
  `location` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '改造地点',
  `content` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '改造内容',
  `is_delay` varchar(1) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT 'N' COMMENT '是否延期',
  `delay_email` varchar(1) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT 'N' COMMENT '是否发送延期通知邮件',
  `director` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '责任人',
  `delay_day` int(1) NULL DEFAULT NULL COMMENT '延期天数，范围：1-7',
  `status` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT 'wait' COMMENT '单据状态：wait待完成、complete完工、withdraw撤单',
  `withdraw_reason` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '撤单原因',
  `is_deleted` varchar(1) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT 'N' COMMENT '是否删除',
  `gmt_created` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '创建时间',
  `gmt_updated` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) ON UPDATE CURRENT_TIMESTAMP(0) COMMENT '修改时间',
  `created_by` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '创建人',
  `updated_by` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '修改人',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 32 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '模具改造' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for mold_repair
-- ----------------------------
DROP TABLE IF EXISTS `mold_repair`;
CREATE TABLE `mold_repair`  (
  `id` bigint(64) NOT NULL AUTO_INCREMENT,
  `code` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '维修编号',
  `line_level` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '线别',
  `fault_desc` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '故障描述',
  `mold_code` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '模具编号',
  `repair_content` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '维修内容',
  `report_time` datetime(0) NOT NULL COMMENT '报告时间',
  `arrive_time` datetime(0) NOT NULL COMMENT '到达时间',
  `finish_time` datetime(0) NOT NULL COMMENT '完成时间',
  `operator` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '操作者',
  `repairtor` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '维修者',
  `auditor_product` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '审核-生产组长',
  `auditor_engine` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '审核-维修工程师',
  `improve_content` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '永久改善内容',
  `lockor` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '锁定人',
  `unlockor` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '解锁人',
  `confirmor` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '确定人',
  `repair_level` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '维修等级',
  `is_deleted` varchar(1) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT 'N' COMMENT '是否删除',
  `created_by` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '创建人',
  `updated_by` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '修改人',
  `gmt_created` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '创建时间',
  `gmt_updated` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) ON UPDATE CURRENT_TIMESTAMP(0) COMMENT '修改时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 20 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '模具维修' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for mold_replace_spare_rel
-- ----------------------------
DROP TABLE IF EXISTS `mold_replace_spare_rel`;
CREATE TABLE `mold_replace_spare_rel`  (
  `id` bigint(64) NOT NULL AUTO_INCREMENT,
  `repair_id` bigint(64) NULL DEFAULT NULL COMMENT '维修ID',
  `maintenance_task_id` bigint(64) NULL DEFAULT NULL COMMENT '保养任务ID',
  `mold_id` bigint(64) NULL DEFAULT NULL COMMENT '模具ID',
  `type` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '更换类型',
  `spare_code` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '备件编号',
  `count` int(11) NOT NULL COMMENT '备件数量',
  `remark` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '备注',
  `is_deleted` varchar(1) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT 'N' COMMENT '是否删除',
  `created_by` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '创建人',
  `updated_by` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '修改人',
  `gmt_created` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '创建时间',
  `gmt_updated` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) ON UPDATE CURRENT_TIMESTAMP(0) COMMENT '修改时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 71 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '模具维修备件更换' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for mqtt_history
-- ----------------------------
DROP TABLE IF EXISTS `mqtt_history`;
CREATE TABLE `mqtt_history`  (
  `id` bigint(64) NOT NULL AUTO_INCREMENT,
  `data` varchar(1024) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT 'mqtt生产数据',
  `gmt_created` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '创建时间',
  `gmt_updated` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) ON UPDATE CURRENT_TIMESTAMP(0) COMMENT '修改时间',
  `created_by` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '创建人',
  `updated_by` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '修改人',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 33 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '登录信息' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for properties
-- ----------------------------
DROP TABLE IF EXISTS `properties`;
CREATE TABLE `properties`  (
  `id` bigint(64) NOT NULL AUTO_INCREMENT,
  `key` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT 'key',
  `value` varchar(2000) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT 'value',
  `is_deleted` varchar(1) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT 'N' COMMENT '是否删除',
  `gmt_created` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '创建时间',
  `gmt_updated` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) ON UPDATE CURRENT_TIMESTAMP(0) COMMENT '修改时间',
  `created_by` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '创建人',
  `updated_by` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '修改人',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 5 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '配置' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for spare_info
-- ----------------------------
DROP TABLE IF EXISTS `spare_info`;
CREATE TABLE `spare_info`  (
  `id` bigint(64) NOT NULL AUTO_INCREMENT,
  `code` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '备件编号',
  `name` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '备件名称',
  `flavor` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '备件规格',
  `material` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '材质',
  `is_deleted` varchar(1) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT 'N' COMMENT '是否删除',
  `gmt_created` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '创建时间',
  `gmt_updated` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) ON UPDATE CURRENT_TIMESTAMP(0) COMMENT '修改时间',
  `created_by` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '创建人',
  `updated_by` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '修改人',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 33 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '备件' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for spare_request
-- ----------------------------
DROP TABLE IF EXISTS `spare_request`;
CREATE TABLE `spare_request`  (
  `id` bigint(64) NOT NULL AUTO_INCREMENT,
  `spare_code` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '备件编号',
  `type` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '出入库类型',
  `location` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '出入库位置',
  `count` int(11) NOT NULL COMMENT '出入库数量',
  `inbound_time` bigint(64) NULL DEFAULT NULL COMMENT '入库时间，表示出的哪一笔记录',
  `remark` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '备注',
  `is_deleted` varchar(1) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT 'N' COMMENT '是否删除',
  `gmt_created` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '创建时间',
  `gmt_updated` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) ON UPDATE CURRENT_TIMESTAMP(0) COMMENT '修改时间',
  `created_by` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '创建人',
  `updated_by` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '修改人',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 88 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '备件库存' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for user_authority
-- ----------------------------
DROP TABLE IF EXISTS `user_authority`;
CREATE TABLE `user_authority`  (
  `id` bigint(64) NOT NULL AUTO_INCREMENT,
  `code` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '编码',
  `key` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '资源标识，唯一',
  `name` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '资源名称',
  `uri` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '资源url',
  `method` varchar(10) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '请求方式',
  `group_name` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '资源组名称',
  `display` varchar(1) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT 'Y' COMMENT '是否返回前端显示',
  `is_deleted` varchar(1) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT 'N' COMMENT '是否删除',
  `gmt_created` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '创建时间',
  `gmt_updated` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) ON UPDATE CURRENT_TIMESTAMP(0) COMMENT '修改时间',
  `created_by` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '创建人',
  `updated_by` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '修改人',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 73 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '用户权限' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for user_authority_rel
-- ----------------------------
DROP TABLE IF EXISTS `user_authority_rel`;
CREATE TABLE `user_authority_rel`  (
  `id` bigint(64) NOT NULL AUTO_INCREMENT,
  `login_name` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '用户登录名',
  `authority_code` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '权限code',
  `is_deleted` varchar(1) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT 'N' COMMENT '是否删除',
  `gmt_created` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '创建时间',
  `gmt_updated` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) ON UPDATE CURRENT_TIMESTAMP(0) COMMENT '修改时间',
  `created_by` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '创建人',
  `updated_by` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '修改人',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 247 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '用户权限关联' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for user_info
-- ----------------------------
DROP TABLE IF EXISTS `user_info`;
CREATE TABLE `user_info`  (
  `id` bigint(64) NOT NULL AUTO_INCREMENT,
  `login_name` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '工号',
  `name` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '姓名',
  `department` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '部门',
  `is_locked` varchar(1) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT 'N' COMMENT '是否锁定',
  `unlock_time` datetime(0) NULL DEFAULT NULL COMMENT '解锁时间',
  `password` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '密码，rsa加密',
  `is_deleted` varchar(1) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT 'N' COMMENT '是否删除',
  `gmt_created` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '创建时间',
  `gmt_updated` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) ON UPDATE CURRENT_TIMESTAMP(0) COMMENT '修改时间',
  `created_by` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '创建人',
  `updated_by` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '修改人',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 22 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '用户' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for mold_inout_bound_task_flow
-- ----------------------------
DROP TABLE IF EXISTS `mold_inout_bound_task_flow`;
CREATE TABLE `mold_inout_bound_task_flow`  (
    `id` bigint(64) NOT NULL AUTO_INCREMENT,
    `job_id` bigint(64) NULL DEFAULT NULL COMMENT '任务ID',
    `task_no` varchar(255) NULL DEFAULT NULL COMMENT '任务编码',
    `flow_order` int(11) NULL DEFAULT NULL COMMENT '任务流转的顺序',
    `src_location` varchar(50) NULL DEFAULT NULL,
    `dest_location` varchar(50) NULL DEFAULT NULL,
    `type` varchar(50) NULL DEFAULT NULL COMMENT '/',
    `agv` tinyint(1) NULL DEFAULT NULL COMMENT 'agv',
    `rfid` varchar(64) NULL DEFAULT NULL COMMENT 'RFID信息',
    `mold_code` varchar(50) NULL DEFAULT NULL COMMENT '模具code',
    `station_unique_code` varchar(5) NULL DEFAULT NULL COMMENT '出入站口唯一标识',
    `error_code` varchar(10) DEFAULT NULL,
    `error_msg` varchar(255) DEFAULT NULL,
    `flag` tinyint(1) NULL DEFAULT 0 COMMENT '流转处理标志',
    `gmt_created` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `gmt_updated` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='任务流转表';

-- ----------------------------
-- Table structure for mold_maintenance_plan_rel
-- ----------------------------
DROP TABLE IF EXISTS `mold_maintenance_plan_rel`;
CREATE TABLE `mold_maintenance_plan_rel`  (
  `id` bigint(64) NOT NULL AUTO_INCREMENT,
  `mold_maintenance_plan_id` bigint(64) NOT NULL COMMENT '模具保养计划ID',
  `mold_id` bigint(64) NOT NULL COMMENT '模具ID',
  `is_deleted` varchar(1) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT 'N' COMMENT '是否删除',
  `created_by` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '创建人',
  `updated_by` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '修改人',
  `gmt_created` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '创建时间',
  `gmt_updated` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) ON UPDATE CURRENT_TIMESTAMP(0) COMMENT '修改时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='模具保养计划关联模具表';

-- ----------------------------
-- Table structure for mold_quality
-- ----------------------------
DROP TABLE IF EXISTS `mold_quality`;
CREATE TABLE `mold_quality`  (
  `id` bigint(64) NOT NULL AUTO_INCREMENT,
  `code` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '质量编号',
  `line_level` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '使用线别',
  `mold_id` bigint(64) NOT NULL COMMENT '模具ID',
  `order_code` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '工单号',
  `production_count` int(11) NOT NULL DEFAULT 0 COMMENT '生产数量',
  `defect_count` int(11) NOT NULL DEFAULT 0 COMMENT '不良数量',
  `defect_desc` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '不良描述',
  `quality_content` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '质量内容',
  `is_deleted` varchar(1) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT 'N' COMMENT '是否删除',
  `created_by` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '创建人',
  `updated_by` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '修改人',
  `gmt_created` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '创建时间',
  `gmt_updated` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) ON UPDATE CURRENT_TIMESTAMP(0) COMMENT '修改时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='模具质量表';

DROP TABLE IF EXISTS `line_product_resume`;
CREATE TABLE `line_product_resume`  (
  `id` bigint(64) NOT NULL AUTO_INCREMENT,
  `line_level` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '线别',
  `shift_id` int(11) NOT NULL COMMENT '班次ID即工单号',
  `shift_min` int(11) NOT NULL COMMENT '班次时长（分钟）',
  `start_time` datetime(0) NOT NULL COMMENT '班次开始时间',
  `end_time` datetime(0) NOT NULL COMMENT '班次结束时间',
  `shift_no` smallint NOT NULL COMMENT '班次号；1-早班；2-晚班',
  `shift_date` varchar(10) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '班次日期',
  `part_code` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '零件名称',
  `qty_ok` int(11) NOT NULL COMMENT '生产数量',
  `qty_nok` int(11) NOT NULL DEFAULT 0 COMMENT '不良数量',
  `is_deleted` varchar(1) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT 'N' COMMENT '是否删除',
  `gmt_created` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '创建时间',
  `gmt_updated` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) ON UPDATE CURRENT_TIMESTAMP(0) COMMENT '修改时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX idx_is_deleted (is_deleted),
  INDEX idx_shift_date (shift_date),
  INDEX idx_shift_id (shift_id),
  INDEX idx_line_level (line_level)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='产线生产履历表';

SET FOREIGN_KEY_CHECKS = 1;
