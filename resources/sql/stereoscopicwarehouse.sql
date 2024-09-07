-- 立库的初始化SQL

-- 1、站点信息

-- 需要替换code中的todo占位符为物理地址编码(待恺博提供)

INSERT INTO `crf_mold`.`mold_stereoscopic_warehouse_station`(`id`, `layer`, `col`, `cell`, `offset`, `name`, `code`, `type`, `outbound_offset`, `inbound_offset`, `stereoscopic_warehouse_id`, `category`, `rfid_device_id`) VALUES (1, 0, 0, 2, 8, '小库小模具出站', 'todo1', 'outbound', 110, NULL, 1, 'B', NULL);
INSERT INTO `crf_mold`.`mold_stereoscopic_warehouse_station`(`id`, `layer`, `col`, `cell`, `offset`, `name`, `code`, `type`, `outbound_offset`, `inbound_offset`, `stereoscopic_warehouse_id`, `category`, `rfid_device_id`) VALUES (2, 0, 0, 1, 10, '小库小模具入站', 'todo2', 'inbound', NULL, 112, 1, 'B', NULL);
INSERT INTO `crf_mold`.`mold_stereoscopic_warehouse_station`(`id`, `layer`, `col`, `cell`, `offset`, `name`, `code`, `type`, `outbound_offset`, `inbound_offset`, `stereoscopic_warehouse_id`, `category`, `rfid_device_id`) VALUES (3, 0, 0, 2, 8, '大库小模具出站', 'todo3', 'outbound', 110, NULL, 2, 'B', NULL);
INSERT INTO `crf_mold`.`mold_stereoscopic_warehouse_station`(`id`, `layer`, `col`, `cell`, `offset`, `name`, `code`, `type`, `outbound_offset`, `inbound_offset`, `stereoscopic_warehouse_id`, `category`, `rfid_device_id`) VALUES (4, 0, 0, 1, 10, '大库小模具入站', 'todo4', 'inbound', NULL, 112, 2, 'B', NULL);
INSERT INTO `crf_mold`.`mold_stereoscopic_warehouse_station`(`id`, `layer`, `col`, `cell`, `offset`, `name`, `code`, `type`, `outbound_offset`, `inbound_offset`, `stereoscopic_warehouse_id`, `category`, `rfid_device_id`) VALUES (5, 0, 0, 4, 12, '大库大模具出站', 'todo5', 'outbound', 114, NULL, 2, 'A', NULL);
INSERT INTO `crf_mold`.`mold_stereoscopic_warehouse_station`(`id`, `layer`, `col`, `cell`, `offset`, `name`, `code`, `type`, `outbound_offset`, `inbound_offset`, `stereoscopic_warehouse_id`, `category`, `rfid_device_id`) VALUES (6, 0, 0, 3, 14, '大库大模具入站', 'todo6', 'inbound', NULL, 115, 2, 'A', NULL);


-- 2、储位信息
-- 生成脚本如下(groovy),需要替换todo占位符,同上：

String template = "INSERT INTO `crf_mold`.`mold_stereoscopic_warehouse_location`(`id`, `layer`, `col`, `cell`, `inbound_station`, `outbound_station`, `mold_id`, `mold_code`, `mold_steroscopic_warehouse_name`, `mold_steroscopic_warehouse_id`, `update_time`, `exist`, `category`) " +
        "VALUES (default, %d, %d, %d, '%s', '%s', NULL, NULL, '%s', %d, now(), 0, '%s');"


// 小库:
println("-- 小库")
String smallInboundCode = "todo2"
String smallOutboundCode = "todo1"
String name = "小库"
category = 'B'
int id = 1
for(int col = 1;col <= 2;col++) {
    for(int cell = 1;cell <= 16;cell++){
        for(int layer = 1;layer <= 15;layer++) {
            // 填充格列层和出入战
            println(String.format(template,layer,col,cell,smallInboundCode,smallOutboundCode,name,id,category))
        }
    }
}

println("-- 大库小模具")
String bigBInboundCode = "todo4"
String bigBoutboundCode = "todo3"
name = "大库"
id = 2
for(int col = 1;col <= 2;col++) {
    for(int cell = 1;cell <= 8;cell++){
        for(int layer = 1;layer <= 15;layer++) {
            // 填充格列层和出入战
            println(String.format(template,layer,col,cell,bigBInboundCode,bigBoutboundCode,name,id,category))
        }
    }
}

println("-- 大库大模具")
category = 'A'
String bigAInboundCode = "todo6"
String bigAoutboundCode = "todo5"
for(int col = 1;col <= 2;col++) {
    for(int cell = 9;cell <= 14;cell++){
        for(int layer = 1;layer <= 10;layer++) {
            // 填充格列层和出入战
            println(String.format(template,layer,col,cell,bigAInboundCode,bigAoutboundCode,name,id,category))
        }
    }
}