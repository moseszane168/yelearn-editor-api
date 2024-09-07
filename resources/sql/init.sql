-- 前端可配置的权限
-- 数据首页
INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
                            (DEFAULT, '010104', 'viewLedger', '查看', '/v1/mold', 'GET', '数据首页', 'Y', 'N', now(), now(), 'admin', 'admin');

-- 出库管理
INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
                            (DEFAULT, '020104', 'moldOutboundView', '查看', '/v1/outbound', 'GET', '出库管理', 'Y', 'N', now(), now(), 'admin', 'admin');
INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
                            (DEFAULT, '020105', 'moldOutbound', '出库权限', '/v1/outbound', 'POST', '出库管理', 'Y', 'N', now(), now(), 'admin', 'admin');
INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
                            (DEFAULT, '020106', 'moldOutbound', '终止出库', '/v1/mold/inoutbound/stop', 'POST', '出库管理', 'Y', 'N', now(), now(), 'admin', 'admin');
INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
                            (DEFAULT, '020107', 'moldOutbound', '取消出库', '/v1/mold/inoutbound/cancel', 'POST', '出库管理', 'Y', 'N', now(), now(), 'admin', 'admin');

-- 入库管理
INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
                            (DEFAULT, '030104', 'moldInboundView', '查看', '/v1/outbound', 'GET', '入库管理', 'Y', 'N', now(), now(), 'admin', 'admin');
INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
                            (DEFAULT, '030105', 'moldInbound', '入库权限', '/v1/outbound', 'POST', '入库管理', 'Y', 'N', now(), now(), 'admin', 'admin');
INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
                            (DEFAULT, '030106', 'moldOutbound', '终止入库', '/v1/mold/inoutbound/stop', 'POST', '入库管理', 'Y', 'N', now(), now(), 'admin', 'admin');
INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
                            (DEFAULT, '030107', 'moldOutbound', '取消入库', '/v1/mold/inoutbound/cancel', 'POST', '入库管理', 'Y', 'N', now(), now(), 'admin', 'admin');

-- 保养管理-待保养
INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
                            (DEFAULT, '040501', 'addPendingMaintainTask', '新增', '/v1/maintenance', 'POST', '保养管理-待保养', 'Y', 'N', now(), now(), 'admin', 'admin');
INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
                            (DEFAULT, '040502', 'updatePendingMaintainTask', '修改', '/v1/maintenance', 'PUT', '保养管理-待保养', 'Y', 'N', now(), now(), 'admin', 'admin');
INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
                            (DEFAULT, '040503', 'deletePendingMaintainTask', '删除', '/v1/maintenance/task/delete', 'POST', '保养管理-待保养', 'Y', 'N', now(), now(), 'admin', 'admin');
INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
                            (DEFAULT, '040504', 'viewPendingMaintainTask', '查看', '/v1/maintenance', 'GET', '保养管理-待保养', 'Y', 'N', now(), now(), 'admin', 'admin');
INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
                            (DEFAULT, '040505', 'togglePendingMaintainTask', '挂起/继续', '/v1/maintenance/task/status', 'POST', '保养管理-待保养', 'Y', 'N', now(), now(), 'admin', 'admin');
INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
                            (DEFAULT, '040506', 'chargePendingMaintainTask', '任务指派', '/v1/maintenance/task/charge', 'POST', '保养管理-待保养', 'Y', 'N', now(), now(), 'admin', 'admin');

-- 保养管理-已保养
INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
                            (DEFAULT, '040602', 'updateDoneMaintainTask', '修改', '/v1/maintenance', 'PUT', '保养管理-已保养', 'Y', 'N', now(), now(), 'admin', 'admin');
INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
                            (DEFAULT, '040603', 'deleteDoneMaintainTask', '删除', '/v1/maintenance/task/delete', 'POST', '保养管理-已保养', 'Y', 'N', now(), now(), 'admin', 'admin');
INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
                            (DEFAULT, '040604', 'viewDoneMaintainTask', '查看', '/v1/maintenance', 'GET', '保养管理-已保养', 'Y', 'N', now(), now(), 'admin', 'admin');
INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
                            (DEFAULT, '040605', 'downloadTaskPDF', '下载PDF', '/v1/maintenance/task/download/pdf', 'POST', '保养管理-已保养', 'Y', 'N', now(), now(), 'admin', 'admin');

-- 保养管理-保养计划
INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
                            (DEFAULT, '040701', 'addMaintainPlan', '新增', '/v1/maintenance/plan', 'POST', '保养管理-保养计划', 'Y', 'N', now(), now(), 'admin', 'admin');
INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
                            (DEFAULT, '040702', 'updateMaintainPlan', '修改', '/v1/maintenance/plan', 'PUT', '保养管理-保养计划', 'Y', 'N', now(), now(), 'admin', 'admin');
INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
                            (DEFAULT, '040703', 'deleteMaintainPlan', '删除', '/v1/maintenance/plan', 'DELETE', '保养管理-保养计划', 'Y', 'N', now(), now(), 'admin', 'admin');
INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
                            (DEFAULT, '040704', 'viewMaintainPlan', '查看', '/v1/maintenance/plan', 'GET', '保养管理-保养计划', 'Y', 'N', now(), now(), 'admin', 'admin');
INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
                            (DEFAULT, '040705', 'toggleMaintainPlan', '暂停/启动', '/v1/maintenance/plan/status', 'POST', '保养管理-保养计划', 'Y', 'N', now(), now(), 'admin', 'admin');

-- 保养管理-保养标准
INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
                            (DEFAULT, '040801', 'addMaintainStandard', '新增', '/v1/maintenance/standard', 'POST', '保养管理-保养标准', 'Y', 'N', now(), now(), 'admin', 'admin');
INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
                            (DEFAULT, '040802', 'updateMaintainStandard', '修改', '/v1/maintenance/standard', 'PUT', '保养管理-保养标准', 'Y', 'N', now(), now(), 'admin', 'admin');
INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
                            (DEFAULT, '040803', 'deleteMaintainStandard', '删除', '/v1/maintenance/standard', 'DELETE', '保养管理-保养标准', 'Y', 'N', now(), now(), 'admin', 'admin');
INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
                            (DEFAULT, '040804', 'viewMaintainStandard', '查看', '/v1/maintenance/standard', 'GET', '保养管理-保养标准', 'Y', 'N', now(), now(), 'admin', 'admin');

-- -- 保养管理-公共
-- INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
--                             (DEFAULT, '040001', 'viewPendingMaintainTask', '查看待保养', '/v1/maintain', 'GET', '保养管理', 'Y', 'N', now(), now(), 'admin', 'admin');
-- INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
--                             (DEFAULT, '040002', 'viewDoneMaintainTask', '查看已保养', '/v1/maintain', 'GET', '保养管理', 'Y', 'N', now(), now(), 'admin', 'admin');
-- INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
--                             (DEFAULT, '040003', 'viewMaintainPlan', '查看保养计划', '/v1/maintain', 'GET', '保养管理', 'Y', 'N', now(), now(), 'admin', 'admin');
-- INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
--                             (DEFAULT, '040004', 'viewMaintainStandard', '查看保养标准', '/v1/maintain', 'GET', '保养管理', 'Y', 'N', now(), now(), 'admin', 'admin');
--
-- -- 保养管理-公共
-- INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
--                             (DEFAULT, '040101', 'addMaintainTask', '新增保养任务', '/v1/maintain', 'POST', '保养管理', 'Y', 'N', now(), now(), 'admin', 'admin');
-- INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
--                             (DEFAULT, '040102', 'updateMaintain', '修改', '/v1/maintain', 'PUT', '保养管理', 'Y', 'N', now(), now(), 'admin', 'admin');
-- INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
--                             (DEFAULT, '040103', 'deleteMaintain', '删除计划/标准', '/v1/maintain', 'DELETE', '保养管理', 'Y', 'N', now(), now(), 'admin', 'admin');
-- INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
--                             (DEFAULT, '040104', 'viewMaintain', '查看', '/v1/maintain', 'GET', '保养管理', 'Y', 'N', now(), now(), 'admin', 'admin');
--
-- -- 保养管理-保养任务
-- INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
--                             (DEFAULT, '040205', 'toggleMaintainTask', '挂起/继续保养任务', '/v1/maintain/task/status', 'POST', '保养管理', 'Y', 'N', now(), now(), 'admin', 'admin');
-- INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
--                             (DEFAULT, '040206', 'chargeMaintainTask', '指派保养任务', '/v1/maintain/task/charge', 'POST', '保养管理', 'Y', 'N', now(), now(), 'admin', 'admin');
-- INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
--                             (DEFAULT, '040207', 'deleteMaintainTask', '删除保养任务', '/v1/maintain/task/delete', 'POST', '保养管理', 'Y', 'N', now(), now(), 'admin', 'admin');
-- INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
--                             (DEFAULT, '040208', 'downloadTaskPDF', '下载已保养任务PDF', '/v1/maintain/task/download/pdf', 'POST', '保养管理', 'Y', 'N', now(), now(), 'admin', 'admin');
--
-- -- 保养管理-保养计划
-- INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
--                             (DEFAULT, '040301', 'addMaintainPlan', '新增保养计划', '/v1/maintain/plan', 'POST', '保养管理', 'Y', 'N', now(), now(), 'admin', 'admin');
-- INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
--                             (DEFAULT, '040305', 'toggleMaintainPlan', '暂停/启动保养计划', '/v1/maintain/plan/status', 'POST', '保养管理', 'Y', 'N', now(), now(), 'admin', 'admin');
--
-- -- 保养管理-保养标准
-- INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
--                             (DEFAULT, '040401', 'addMaintainStandard', '新增保养标准', '/v1/maintain/standard', 'POST', '保养管理', 'Y', 'N', now(), now(), 'admin', 'admin');

-- 维修管理
INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
                            (DEFAULT, '050101', 'addRepair', '新增', '/v1/mold/repair', 'POST', '模具维修', 'Y', 'N', now(), now(), 'admin', 'admin');
INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
                            (DEFAULT, '050102', 'updateRepair', '修改', '/v1/mold/repair', 'PUT', '模具维修', 'Y', 'N', now(), now(), 'admin', 'admin');
INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
                            (DEFAULT, '050103', 'deleteRepair', '删除', '/v1/mold/repair', 'DELETE', '模具维修', 'Y', 'N', now(), now(), 'admin', 'admin');
INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
                            (DEFAULT, '050104', 'viewRepair', '查看', '/v1/mold/repair', 'GET', '模具维修', 'Y', 'N', now(), now(), 'admin', 'admin');

-- 模具立库管理
-- INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
--                             (DEFAULT, '060101', 'addSterescopic', '新增', '/v1/sterescopic/location', 'POST', '模具立库管理', 'Y', 'N', now(), now(), 'admin', 'admin');
INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
                            (DEFAULT, '060102', 'updateSterescopic', '修改', '/v1/sterescopic/location', 'PUT', '模具立库管理', 'Y', 'N', now(), now(), 'admin', 'admin');
-- INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
--                             (DEFAULT, '060103', 'deleteSterescopic', '删除', '/v1/sterescopic/location', 'DELETE', '模具立库管理', 'Y', 'N', now(), now(), 'admin', 'admin');
INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
                            (DEFAULT, '060104', 'viewSterescopic3D', '查看立库动画', '/v1/sterescopic/location', 'GET', '模具立库管理', 'Y', 'N', now(), now(), 'admin', 'admin');
INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
                            (DEFAULT, '060105', 'viewSterescopicLocation', '查看模具仓', '/v1/sterescopic/location', 'POST', '模具立库管理', 'Y', 'N', now(), now(), 'admin', 'admin');

-- 生产履历
INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
                            (DEFAULT, '070104', 'viewProductresume', '查看', '/v1/mold/productresume/page', 'GET', '生产履历', 'Y', 'N', now(), now(), 'admin', 'admin');
INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
                            (DEFAULT, '070105', 'exportProductresume', '导出明细', '/v1/mold/productresume/excel', 'GET', '生产履历', 'Y', 'N', now(), now(), 'admin', 'admin');

-- 台账管理
INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
                            (DEFAULT, '080101', 'addLedger', '新增', '/v1/mold', 'POST', '台账管理', 'Y', 'N', now(), now(), 'admin', 'admin');
INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
                            (DEFAULT, '080102', 'updateLedger', '修改', '/v1/mold', 'PUT', '台账管理', 'Y', 'N', now(), now(), 'admin', 'admin');
INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
                            (DEFAULT, '080103', 'deleteLedger', '删除', '/v1/mold', 'DELETE', '台账管理', 'Y', 'N', now(), now(), 'admin', 'admin');
INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
                            (DEFAULT, '080104', 'viewLedger', '查看', '/v1/mold', 'GET', '台账管理', 'Y', 'N', now(), now(), 'admin', 'admin');
INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
                            (DEFAULT, '080105', 'importLedger', '导入台账', '/v1/mold/excel', 'GET', '台账管理', 'Y', 'N', now(), now(), 'admin', 'admin');
INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
                            (DEFAULT, '080106', 'exportLedger', '导出台账', '/v1/mold/excel', 'POST', '台账管理', 'Y', 'N', now(), now(), 'admin', 'admin');

-- 模具改造
INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
                            (DEFAULT, '090101', 'addRemodel', '新增', '/v1/mold/remodel', 'POST', '模具改造', 'Y', 'N', now(), now(), 'admin', 'admin');
INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
                            (DEFAULT, '090102', 'updateRemodel', '修改', '/v1/mold/remodel', 'PUT', '模具改造', 'Y', 'N', now(), now(), 'admin', 'admin');
INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
                            (DEFAULT, '090103', 'deleteRemodel', '删除', '/v1/mold/remodel', 'DELETE', '模具改造', 'Y', 'N', now(), now(), 'admin', 'admin');
INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
                            (DEFAULT, '090104', 'viewRemodel', '查看', '/v1/mold/remodel/page', 'GET', '模具改造', 'Y', 'N', now(), now(), 'admin', 'admin');
INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
                            (DEFAULT, '090105', 'completeRemodel', '完工改造', '/v1/mold/remodel/completed', 'POST', '模具改造', 'Y', 'N', now(), now(), 'admin', 'admin');
INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
                            (DEFAULT, '090106', 'withdrawRemodel', '撤单改造', '/v1/mold/remodel/withdraw', 'POST', '模具改造', 'Y', 'N', now(), now(), 'admin', 'admin');

-- 知识库管理
INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
                            (DEFAULT, '100101', 'addLibrary', '新增', '/v1/knowledge', 'POST', '知识库管理', 'Y', 'N', now(), now(), 'admin', 'admin');
INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
                            (DEFAULT, '100102', 'updateLibrary', '修改', '/v1/knowledge', 'PUT', '知识库管理', 'Y', 'N', now(), now(), 'admin', 'admin');
INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
                            (DEFAULT, '100103', 'deleteLibrary', '删除', '/v1/knowledge', 'DELETE', '知识库管理', 'Y', 'N', now(), now(), 'admin', 'admin');
INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
                            (DEFAULT, '100104', 'viewLibrary', '查看', '/v1/knowledge', 'GET', '知识库管理', 'Y', 'N', now(), now(), 'admin', 'admin');

-- 备件履历
-- INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
--                             (DEFAULT, '110101', 'addSpare', '新增', '/v1/spare', 'POST', '备件履历', 'Y', 'N', now(), now(), 'admin', 'admin');
-- INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
--                             (DEFAULT, '110102', 'updateSpare', '修改', '/v1/spare', 'PUT', '备件履历', 'Y', 'N', now(), now(), 'admin', 'admin');
-- INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
--                             (DEFAULT, '110103', 'deleteSpare', '删除', '/v1/spare', 'DELETE', '备件履历', 'Y', 'N', now(), now(), 'admin', 'admin');
INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
                            (DEFAULT, '110104', 'viewSpare', '查看', '/v1/spare/resume/page', 'POST', '备件履历', 'Y', 'N', now(), now(), 'admin', 'admin');
INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
                            (DEFAULT, '110105', 'viewSpare', '导出备件', '/v1/spare/resume/excel', 'GET', '备件履历', 'Y', 'N', now(), now(), 'admin', 'admin');

-- 备件管理
-- INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
--                             (DEFAULT, '120101', 'addSpare', '新增', '/v1/spare', 'POST', '备件管理', 'Y', 'N', now(), now(), 'admin', 'admin');
INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
                            (DEFAULT, '120102', 'updateSpare', '修改', '/v1/spare', 'PUT', '备件管理', 'Y', 'N', now(), now(), 'admin', 'admin');
INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
                            (DEFAULT, '120103', 'deleteSpare', '删除', '/v1/spare', 'DELETE', '备件管理', 'Y', 'N', now(), now(), 'admin', 'admin');
INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
                            (DEFAULT, '120104', 'viewSpare', '查看', '/v1/spare/', 'GET', '备件管理', 'Y', 'N', now(), now(), 'admin', 'admin');
INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
                            (DEFAULT, '120201', 'addSpareRequest', '新增库存清单', '/v1/spare/', 'POST', '备件管理', 'Y', 'N', now(), now(), 'admin', 'admin');
INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
                            (DEFAULT, '120301', 'addSpare', '新增基础资料', '/v1/spare/', 'POST', '备件管理', 'Y', 'N', now(), now(), 'admin', 'admin');

-- 配置中心
INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
                            (DEFAULT, '130101', 'addUser', '添加', '/v1/dict', 'POST', '配置中心', 'Y', 'N', now(), now(), 'admin', 'admin');
INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
                            (DEFAULT, '130102', 'updateUser', '编辑', '/v1/dict', 'PUT', '配置中心', 'Y', 'N', now(), now(), 'admin', 'admin');
INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
                            (DEFAULT, '130103', 'deleteUser', '删除', '/v1/dict', 'DELETE', '配置中心', 'Y', 'N', now(), now(), 'admin', 'admin');
INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
                            (DEFAULT, '130104', 'viewUser', '查看', '/v1/dict', 'GET', '配置中心', 'Y', 'N', now(), now(), 'admin', 'admin');

-- 前端不可配置，管理员专用的权限
-- 用户管理
INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
                            (DEFAULT, '140101', 'addUser', '添加', '/v1/user', 'POST', '用户管理', 'Y', 'N', now(), now(), 'admin', 'admin');
INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
                            (DEFAULT, '140102', 'updateUser', '编辑', '/v1/user', 'PUT', '用户管理', 'Y', 'N', now(), now(), 'admin', 'admin');
INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
                            (DEFAULT, '140103', 'deleteUser', '删除', '/v1/user', 'DELETE', '用户管理', 'Y', 'N', now(), now(), 'admin', 'admin');
INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
                            (DEFAULT, '140104', 'viewUser', '查看', '/v1/user', 'GET', '用户管理', 'Y', 'N', now(), now(), 'admin', 'admin');
INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
                            (DEFAULT, '140105', 'updateUserAuthority', '权限配置', '/v1/user/authority', 'PUT', '用户管理', 'Y', 'N', now(), now(), 'admin', 'admin');
INSERT INTO `user_authority`(`id`, `code`, `key`, `name`, `uri`, `method`, `group_name`, `display`,  `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
                            (DEFAULT, '140106', 'resetPwd', '重置密码', '/v1/user/pwd', 'PUT', '用户管理', 'Y', 'N', now(), now(), 'admin', 'admin');

-- 公私钥
INSERT INTO `properties`(`id`, `key`, `value`, `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
                        (default, 'rsa_public_key', 'MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAvrKWFpLztzPQ3UoKTt1Q\\nNAxfSoiqfqTefq0w92Jj9lYkVC6fQR2NBdQKsZeT6q8jAcEYWVnXHiGT3BRN+WLE\\ngMU52uSJSfdTulj9Nt3yYOtWzCmM2xGri14jhIa3IdnhfcX2XNoNB/fNgEeo0o+t\\n2+9yHSe9wW6RtPCBE8U46MGwTnFZtlBXikA66135ubyYtYd2LrPN7nPYEkNL8TtY\\nMbMTDmyzgUa5lOXypcbc+zd+QOgbIXtCkNvN17Nk9tFOB5yPN3toLNzOkmMhYREv\\n0dXndJ8klDynrPkkCfDxdIvq9RjPNoqTxhweCeXHDAekIn3lg54VC3CacQrPluSu\\npwIDAQAB', 'N', now(), now(), 'admin', 'admin');
INSERT INTO `properties`(`id`, `key`, `value`, `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
                        (default, 'rsa_private_key', 'MIIEuwIBADALBgkqhkiG9w0BAQEEggSnMIIEowIBAAKCAQEAvrKWFpLztzPQ3UoK\\nTt1QNAxfSoiqfqTefq0w92Jj9lYkVC6fQR2NBdQKsZeT6q8jAcEYWVnXHiGT3BRN\\n+WLEgMU52uSJSfdTulj9Nt3yYOtWzCmM2xGri14jhIa3IdnhfcX2XNoNB/fNgEeo\\n0o+t2+9yHSe9wW6RtPCBE8U46MGwTnFZtlBXikA66135ubyYtYd2LrPN7nPYEkNL\\n8TtYMbMTDmyzgUa5lOXypcbc+zd+QOgbIXtCkNvN17Nk9tFOB5yPN3toLNzOkmMh\\nYREv0dXndJ8klDynrPkkCfDxdIvq9RjPNoqTxhweCeXHDAekIn3lg54VC3CacQrP\\nluSupwIDAQABAoIBAHdNsRp0W2ctUqlvDd3jFa9KYj92GvxaVxx3a+AJPTK7F8VW\\n2alaPIT98KbEhvTXFxac4Ifd7fha129jgJjaEsfhG933BnEw+7/ktp4h4uaBtW7L\\nO+U+O81YWu4pfd7+udT/Ca9zd52ZiYaMznDVFNc5CXJ2D4A5lYzWvlpJE96BYj5V\\nu6SeSJCQtxv8Y4Ey7n/Er6HkjMwWCoWiaClpPzJeftJUibHxH063HLPvKGXZFGX5\\n9vF04xyfzuJQiLgbCTXrwwTaxYO+8UnS7esqlge4fkpzhYhFtTnyC3NSEgj91tJH\\nYeoeKFsbgPVdu3Evk4u0Gxdna0jEGnDqf4CZCcECgYEA3OpmD0edOioJrTqURy5G\\nqu1pKd+67QhmvtDITV7wKQfpnJEk/p4LEv2WfgEMkIouj14H3EJiuECexKlve0uq\\nw9wdofcpBCxCJkotmL/KSd6I6eALKa+XwK6YgtyEmrggixaicZfzpmNfaMAShOIZ\\npiSIQWQNgfIsBvoBZS+qDkcCgYEA3Pui4gWxVyNA7Sc0KBRm8IGAGQGXnYSzGotE\\n/+8JKk7YdYqaHiEY4+ZCja484nKWK+ZZ+Wg434TJTGzWaN/uNOY3H78yNttYRgnB\\n/QzkHVlp/IgiqyfOge86elH+7wCAh4wbJK+UDvaCRCifdiS9sTLKQsR/GBOyLwvw\\nKDTWrKECgYAw4XiFpv3mEckkWFLY0Sd3yKI9TrDIo9RAImg/nmMbYRHSv9bks8mV\\ngSDcbpT+ImUc+dxZYyL+y+WVdDwjluGJBtpTrSGZN8XHPSCLrNwwrhmzTgyKQ70b\\nOEaspeh9Z4Jj5DU7VzjlNxW0UtOGLZUpSuoPNfk7KH+PZ6AJaJuDHwKBgQCd3jsP\\n42c8zBefFInDNEgR+0HrG2MYCev1w5bIjBjtG6Sx3BGcAqMIdMAI/XfLgnbb59VR\\nQu6WaANy0LIf/BHtwqWQzYNvAyY96symHeZ9PRplaU/zHB4AX0pUhm1sitxHeYUO\\noUxRoDORw7+fpEHL7G/oYP420iNSTuIDpzPR4QKBgHUAONjv7DGVy5QbMKWdpkMk\\n9v/KcvtEhpfJSiSJYkK7Kbjnb5thAmK/b6LxU/jV8hf4qR9i5qsTLe4VE3LeHxmH\\nS7ra8PtGKRX4AZ/2i43+XZP3YqjSPPNR0LNLFY3HAmPYhFoWkLsLJl1Wzm6GteAk\\n3OgA0tYDb7foAq8RProS', 'N', now(), now(), 'admin', 'admin');

-- 超时邮件
INSERT INTO `properties`(`id`, `key`, `value`, `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
                        (DEFAULT, 'maintenance_timeout_email', '', 'N', now(), now(), 'admin', 'admin');
INSERT INTO `properties`(`id`, `key`, `value`, `is_deleted`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
                        (DEFAULT, 'remodel_timeout_email', '', 'N', now(), now(), 'admin', 'admin');

-- 默认账号admin,密码123456;和上面公私钥绑定
INSERT INTO `user_info`(`id`, `login_name`, `name`, `department`, `password`, `is_deleted`, `is_root`, `gmt_created`, `gmt_updated`, `created_by`, `updated_by`) VALUES
                       (1, 'admin', 'admin', 'IT', 'e10adc3949ba59abbe56e057f20f883e', 'N', 'Y', now(), now(), 'admin', 'admin');