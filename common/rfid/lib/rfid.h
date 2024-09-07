
/**
 * @file rfid.h
 *
 * @brief rfid reader sdk header file
 */

/**@mainpage rfid读写器sdk c库/rfid reader sdk c library
 *
 * # 简介/brief introduction
 *
 * 本库提供rfid读写器所有的功能包括：<p>
 *
 * This library provides all functions of RFID reader including: <p>
 *
 * - rfid标签读、写、盘点、锁定、毁灭等功能
 * - rfid读写器管理功能
 * - 一些工具函数
 *<p>
 * - RFID tag reading, writing, inventory, locking, destruction and other functions
 * - RFID reader management function
 * - Some utility functions
 *
 * # 基本用法/Basic Usage
 *
 * ## 设备扫描功能/Device scanning
 *
 * 可以扫描到局域网内的已联网的设备<p>
 * Can scan to the network of devices in the LAN<p>
 *
 * - 定义回调/Define the callback
 *
 * @code
 static void test_device_scan_callback(rfid_device device) {
  printf("device: %s %d %s\n", device.ip, device.port, device.id);
  free_rfid_device(device);
}
 * @endcode
 *
 * > 注意释放内存`free_rfid_device`，避免内存泄漏<p>
 * > Pay attention to release memory 'free_rfid_device' to avoid memory leakage<p>
 *
 * - 开始扫描/Start Scan
 *
 * @code
 int ret = start_device_scan(test_device_scan_callback);
 * @endcode
 *
 * - 结束扫描/Stop Scan
 *
 * @code
 ret = stop_device_scan();
 * @endcode
 *
 * ## 连接设备/Connecting
 *
 * 本库支持tcp，udp，rs232，rs485方式连接设备（实际看设备支持哪种）<p>
 * This library supports TCP, UDP, RS232, RS485 to connect the device (actually depends on which device supports)<p>
 *
 * - tcp
 *
 * @code
 * int ret = connect_net("192.168.1.223", 1969, tcp, &id);
 * @endcode
 *
 * - udp
 *
 * @code
 * int ret = connect_net("192.168.1.223", 1969, udp, &id);
 * @endcode
 *
 * - rs232
 *
 * @code
 * int ret = connect_rs((char *) "COM3", baud115200, 0, &id);
 * @endcode
 *
 * - rs485
 *
 * @code
 * int ret = connect_rs((char *) "COM3", baud115200, 33, &id);
 * @endcode
 *
 * 连接成功之后，返回本次连接的id，该id用作标识本连接<p>
 * After successful connection, return the ID of this connection, which is used to identify this connection
 *
 * ## 标签操作/Tag operation
 *
 * ### 列出附近的标签/List nearby tags
 *
 *  本功能只能在命令模式下使用<p>
 *  This feature can only be used in command mode<p>
 *
 * @code
 unsigned char password[4] = {0, 0, 0, 0};
  int count;
  tag **tags;
  int ret =
      list6c(id, memory_bank_epc, 4, 4, password, 4, &tags, &count);
  // assert(ret == err_ok);
  if (ret != err_ok) {
    return;
  }
  for (int i = 0; i < count; ++i) {
    char epc_hex_str[256];
    to_hex_str(tags[i]->id, tags[i]->len, epc_hex_str);
    printf("%d: list 6c epc: [%s]\n", i, epc_hex_str);
  }
 * @endcode
 *
 * > 注意使用完之后释放内存 `free_tags(tags);` <p>
 * > Free memory after use `free_tags(tags);` <p>
 *
 * ### 读标签/Reading
 *
 * 以下代码展示读取user区0-12字节的数据<p>
 * The following code shows reading 0-12 bytes of data from the user memory<p>
 *
 * @code
  unsigned char user[256];
  ret = read6c(id, memory_bank_user, 0, 12, epc, epc_len,
                 password, 4, user, &dataLen);
 * @endcode
 *
 * ### 写标签/Writing
 *
 *  以下代码展示写入user区0开始12个字节的数据<p>
 *  The following code shows writing 0-12 bytes of data from the user memory<p>
 *
 * @code
  ret = write6c(id, memory_bank_user, 0, epc,epc_len, user, 12,
                  password, 4);
 * @endcode
 *
 * ### 监听自动模式下输出/Listen the output in automatic mode
 *
 * 当读写器处于自动工作模式时，以下代码可以监听自动输出的数据<p>
 * When the reader is in auto mode, the following code can listen for auto-output data
 *
 * > 用户可以自定义输出的数据字段，可以通过`set_output_format`或者读写器管理程序来修改,
 * > 参考后面的章节[自定义上报格式](@ref lable_define_output_format) <p>
 * > Users can customize the output data fields. This can be modified using 'set_output_format' or the reader
 * > hypervisor. Refer to the following section [Custom Output Format](@ref lable_define_output_format)
 *
 * - 定义回调/Define the callback
 *
 * @code
 void testAutoReadCB(auto_read_output item) {
  printf("testAutoReadCB auto read: %s\n", item.epc);
}
 * @endcode
 *
 * - 开始监听/Start Scan
 *
* @code
 ret = start_listen_auto_read(id, testAutoReadCB);
 * @endcode
 *
 * - 结束监听/Stop Scan
 *
 * @code
 * ret = stop_listen_auto_read(id);
 * @endcode
 *
 *
 * ### 自定义上报格式/Custom report format ###                         {#lable_define_output_format}
 *
 * - 设置格式字段/format fields setting
 *
 * @code
 output_format format;
  memset(&format, 0, sizeof(output_format));
  format.format = format.format | format_flag_ant | format_flag_epc |
                  format_flag_ip | format_flag_tid;
 * @endcode
 *
 * - 测试某个字段是否设置了/check format field
 *
 *@code
 assert(format.format & format_flag_ant);
  assert(format.format & format_flag_epc);
  assert(format.format & format_flag_ip);
  assert(format.format & format_flag_tid);
 * @endcode
 *
 *
 * - 设置格式/output format setting
 *
 * @code
  int ret = set_output_format(id, format);
  assert(ret == err_ok);
 * @endcode
 *
 *
 */

#ifndef _RFID_H_
#define _RFID_H_

#if !defined(WIN32)
#define __declspec(dllexport)
#endif

#include <stddef.h>
#include <stdint.h>

#define err_ok 0

#define MAX_EPC_LEN 256
#define MAX_USER_LEN 256
#define MAX_TID_LEN 256

#define MAX_HF_UID_LEN 10

#define MAX_IP_LEN 15
#define MAX_WIFI_PASSWORD_LEN 64
#define MAX_WIFI_SSID_LEN 32
#define MAX_ID_LEN 24
#define MAX_CUSTOM_LEN 24
#define MAX_TRIGGER_COUNT 16
typedef struct _VariableSizedStructure
{
  int length;
  char buffer[1];
} VariableSizedStructure;

/**
 * @brief 传输协议/transport protocol
 */
typedef enum transport_protocol
{
  tcp,
  udp
} transport_protocol;

/**
 * @brief 波特率/Baud rate
 */
typedef enum baud
{
  baud9600,
  baud19200,
  baud38400,
  baud57600,
  baud115200
} baud;

/**
 * @brief 硬件接口/hardware interface
 */
typedef enum port
{
  /// RS232接口/RS232 interface
  rs232,
  /// RS485接口/RS485 interface
  rs485,
  /// 以太网接口/Ethernet network interface
  rj45,
  /// 韦根26/Wiegand-26
  wg26,
  /// 韦根34/Wiegand-34
  wg34,
  /// WiFi接口/WiFi interface
  wifi,
  /// GSM接口(4g/5g)/GSM interface (4g/5g)
  gsm4g,
} port;

/**
 * @brief 继电器状态/relay status
 */
typedef enum relay_status
{
  /// 继电器断开/Relay off
  relay_closed,
  /// 继电器闭合/relay closure
  relay_open
} relay_status;

/**
 * @brief 电平/electrical level
 */
typedef enum electrical_level
{
  /// 高电平/high level
  high,
  /// 低电平/low level
  low
} electrical_level;

/**
 * @brief 自动模式下输出数据/Output data in automatic mode
 *
 *
 */
typedef struct auto_read_output
{
  /// 天线号/Antenna NO
  int ant;
  /**
   * 表示开关触发模式下哪一路被触发了，它的值可能是[1,2,3]。该字段只在开关触发模式下才会上报。<p>
   * Indicates which triggers are triggered in trigger mode. Its value may be [1,2,3]. This field is reported only in
   * triggered mode.
   */
  int fin[MAX_TRIGGER_COUNT];
  /// 被触发的触发器数量/Number of triggers
  int fin_count;
  /**
   * 表示通道门模式下的触发方向，从触发1到触发2为`In`，从触发2到触发1为`Out`。该字段只在通道门模式下才会上报。<p>
   * Indicates the trigger direction In the gate mode. from trigger 1 to trigger 2 is 'IN', and from trigger 2 to
   * trigger 1 is 'OUT'. This field is only reported in Gateway mode
   */
  char door[4];
  /// 读写器的IP地址/The IP address of the reader
  char ip[MAX_IP_LEN + 1];
  /// 高频卡UID
  char uid[16 + 1];
  /// EPC
  char epc[MAX_EPC_LEN * 2 + 1];
  /// TID
  char tid[MAX_TID_LEN * 2 + 1];
  /// USER
  char user[MAX_USER_LEN * 2 + 1];
  /// 读写器的设备ID//The device ID of the reader
  char id[MAX_ID_LEN + 1];
  /// 标签的信号强度/rssi
  int rssi;
  /// 识别到该标签的时间戳(unix time)/timestamp
  int ts;
  /// 标签的类型。1--6C,2--6B/The type of the tag. 1--6C,2--6B
  int tagType;
  /// 用户自定义字段1/User defined field 1
  char custom1[MAX_CUSTOM_LEN + 1];
  /// 用户自定义字段2/User defined field 2
  char custom2[MAX_CUSTOM_LEN + 1];
  /// 用户自定义字段3//User defined field 3
  char custom3[MAX_CUSTOM_LEN + 1];
  /// 用户自定义字段4//User defined field 4
  char custom4[MAX_CUSTOM_LEN + 1];
  /// 用户自定义字段5//User defined field 5
  char custom5[MAX_CUSTOM_LEN + 1];
} auto_read_output;

typedef struct report_option
{
  char ip[64];
  int port;
  int enable;
} report_option;

typedef struct device_config
{
  char name[64];
  int wifi_available;
  int g4_available;
  int wg_available;
  int supernet_available;
  int ant_count;
  int trigger_count;
  int relay_count;
  int di;    // device type
  int proto; // 0-uvp 1-p218 2-r2000
} device_config;

/**
 * 自动模式数据监听callback/Automatic mode data monitoring callback
 */
typedef void (*auto_read_callback)(auto_read_output);

typedef struct testCallback
{
  int (*read)(unsigned char*, size_t);
} testCallback;

/**
 * @brief 表示扫描到的读写器设备/scanned reader device
 *
 * 主要是扫描设备回调中使用<p>
 * It is mainly used for scanning device callbacks
 */
typedef struct rfid_device
{
  /// 设备IP地址/Device IP address
  char* ip;
  /// tcp/udp端口/tcp/udp port
  int port;
  /// 读写器型号/Reader model
  char* model;
  /// 读写器ID/ID
  char* id;
  /// 485地址/rs485 address
  int rs485_address;
  // 232波特率/rs232 baud rate
  int rs232_baud;
  /// 485波特率/rs485 baud rate
  int rs485_baud;
  int ti;
  int proto;
} rfid_device;

/// 设备扫描回调函数/Device scan callback
typedef void (*rfid_device_scan_callback)(rfid_device);

/// 继电器自动控制参数/Relay automatic control parameters
typedef struct relay_option
{
  /**
   * @brief 继电器的用途/Use of Relays
   *
   * - 0不自动控制继电器
   * - 1当读到标签时，吸合继电器,
   * - 2当读到满足报警条件的标签时，吸合继电器<p>
   *
   * - 0 disable
   * - 1 When the label is read, pull in the relay.
   * - 2 When a label that meets the alarm condition is read, pull in the relay
   *
   */
  int tp;
  /// 表示闭合超时时间，单位：秒/the closure timeout time(second)
  int timeout;
} relay_option;

/// 工作模式/work mode
typedef enum work_mode
{
  /// 命令模式/command mode
  command_mode,
  /**
   * 自动模式（连续读卡）：读写器持续地自动读卡。<p>
   * Automatic mode (continuous card reading): The reader will continuously and automatically read the card.
   */
  auto_mode,
  /**
   * 自动模式（开关触发读卡）：当读写器的触发IO满足触发条件时，读一段时间的卡。<p>
   * Automatic mode (reading when triggered) : when the trigger IO of the reader meets the trigger condition,
   * read the card for a period of time.
   */
  trigger_mode,
  /**
   * 自动模式（双向触发读卡）：当读写器的触发IO1被触发时，开始读卡，触发IO2被触发时，停止读卡。或者，读写器的触发IO2被触发时，开始读卡，触发IO1被触发时，停止读卡。触发条件和时间可通过2FH命令设置。<p>
   *
   * Automatic mode (reading when bidirectional trigger is triggered ):When the reader's trigger IO1 is triggered,
   *  reading will be started, and when the trigger IO2 is triggered, reading will be stopped. or, conversely. Trigger
   * conditions and timeout can be set by the 2FH command
   */
  bidirection_trigger_mode,
} work_mode;

/// 频率区域/frequency region
typedef enum freq_region
{
  /// FCC
  region_usa,
  /// ETSI
  region_eu,
  /// CHN
  region_cn,
  /// USER DEFINED
  region_user,
} freq_region;

/**
 * @brief 自定义频率设置方式/Custom frequency Settings mode
 */
typedef enum freq_setting_type
{
  /// 自定义频点：完全自定义 / Custom Frequency Points: Fully custom
  freq_setting_user,
  /// 区域内自定义频点 / Customize frequency points within a region
  freq_setting_in_region,
  /// 区域内自动全频 / Automatic full frequency in the region
  freq_setting_in_region_auto,
} freq_setting_type;

/**
 * @brief 频率设置参数/Frequency setting parameter
 *
 * 根据不同的 @ref freq_setting_type 来实现不同的频率设置方式.<br>
 * Depending on the @ref freq_setting_type, different frequency Settings are implemented.
 *
 * - freq_setting_user
 *
 *  本模式下，region固定为region_user，freq_space,freq_count,start_freq有效<br>
 *  In this mode,region is fixed as region_user. freq_space,freq_count,start_freq are valid.
 *
 * - freq_setting_in_region
 *
 *  本模式下，region，start_freq_seq，end_freq_seq字段有效<br>
 *
 *  In this mode, region，start_freq_seq，end_freq_seq are valid.
 *
 * - freq_setting_in_region_auto
 *
 *  本模式下，仅region字段有效<br>
 *  In this mode, only region，is valid.
 */
typedef struct frequency_option
{
  /// 频段 region
  freq_region region;
  /// 频点间隔 FreqSpace，频点间隔 = FreqSpace * 10KHz。<br>
  ///  frequency points gap (unit is 10KHz)
  int freq_space;
  /**
   * @brief 起始频点/Start frequency point
   *
   * @details
   *
   * FCC范围0-52，ETSI范围0-6，CHN范围0-10， 具体请参考一下对应表<p>
   *
   * FCC Range 0-52, ETSI Range 0-6, CHN Range 0-10, Please refer to the corresponding table for details<p>
   *
   - FCC

    频点/frequency point|频率/frequency
    ------------- | -------------
    0	      |  865.00 MHz
    1	      |  865.50 MHz
    2	      |  866.00 MHz
    3	      |  866.50 MHz
    4	      |  867.00 MHz
    5	      |  867.50 MHz
    6	      |  868.00 MHz

    - ETSI

    频点/frequency point|频率/frequency|频点/frequency point|频率/frequency|频点/frequency point|频率/frequency|
    ------------- | -------| ------------- | -------|------------- | -------|
    0	        |902.00 MHz|20	|912.00 MHz|40	|922.00 MHz|
    1	        |902.50 MHz|21	|912.50 MHz|41	|922.50 MHz|
    2	        |903.00 MHz|22	|913.00 MHz|42	|923.00 MHz|
    3	        |903.50 MHz|23	|913.50 MHz|43	|923.50 MHz|
    4	        |904.00 MHz|24	|914.00 MHz|44	|924.00 MHz|
    5	        |904.50 MHz|25	|914.50 MHz|45	|924.50 MHz|
    6	        |905.00 MHz|26	|915.00 MHz|46	|925.00 MHz|
    7	        |905.50 MHz|27	|915.50 MHz|47	|925.50 MHz|
    8	        |906.00 MHz|28	|916.00 MHz|48	|926.00 MHz|
    9	        |906.50 MHz|29	|916.50 MHz|49	|926.50 MHz|
    10	        |907.00 MHz|30	|917.00 MHz|50	|927.00 MHz|
    11	        |907.50 MHz|31	|917.50 MHz|51	|927.50 MHz|
    12	        |908.00 MHz|32	|918.00 MHz|52	|928.00 MHz|
    13	        |908.50 MHz|33	|918.50 MHz|-   |-         |
    14	        |909.00 MHz|34	|919.00 MHz|-   |-         |
    15	        |909.50 MHz|35	|919.50 MHz|-   |-         |
    16	        |910.00 MHz|36	|920.00 MHz|-   |-         |
    17	        |910.50 MHz|37	|920.50 MHz|-   |-         |
    18	        |911.00 MHz|38	|921.00 MHz|-   |-         |
    19	        |911.05 MHz|39	|921.50 MHz|-   |-         |

    - CHN

       频点/frequency point|频率/frequency
    ------------- | -------------
    0	|920.00 MHz|
    1	|920.50 MHz|
    2	|921.00 MHz|
    3	|921.50 MHz|
    4	|922.00 MHz|
    5	|922.50 MHz|
    6	|923.00 MHz|
    7	|923.50 MHz|
    8	|924.00 MHz|
    9	|924.50 MHz|
    10	|925.00 MHz|
    *
   */
  int start_freq_seq;
  /// 结束频点。FCC范围0-52，ETSI范围0-6，CHN范围0-10。<br>
  /// end frequency point: FCC Range 0-52, ETSI Range 0-6, CHN Range 0-10
  int end_freq_seq;
  /// 频点数量 FreqQuantity
  int freq_count;
  /// 起始频率 StartFreq 单位为 KHz/ start frequency(unit KHz)
  int start_freq;
  /// 设置模式/ setting mode
  freq_setting_type setting_type;
} frequency_option;

/**
 * @brief 选择标签算法选项/ label algorithm option
 */
typedef struct tag_algorithm
{
  /**
   * 是否打开少量标签识别算法，0--否，1--是。仅支持部分模块<br>
   * Whether to open a small number of tag recognition algorithm, 0-- no, 1-- yes. Only some modules are supported
   */
  int less_tag;
  /**
   * 是否打开快速读取EPC＋TID+USER功能，0--否，1--是<br>
   * Whether to enable quick read EPC+TID+USER function, 0-- No, 1-- Yes
   */
  int fast_tid;
  /**
   * 是否打开快速天线切换功能，0--否，1--是<br>
   * Whether to turn on fast antenna switch, 0-- no, 1-- yes
   */
  int fast_ant_switch;
} tag_algorithm;

/**
 * @brief 标签过滤器/Label filter
 */
typedef struct filter
{
  /**
   *  过滤器的用途。
   *
   *  - 0 - 仅输出匹配的标签；
   *  - 1 - 输出所有标签，当识别到匹配的标签时，如果开启了自动控制继电器功能，则会触发继电器工作。<br>
   *
   * Purpose of the filter
   *
   * - 0 - Only matched labels are output
   * - 1 - all labels are output.
   */
  int enable;
  /// 掩码长度/Mask length
  int len;
  /// 掩码地址/ Mask Address
  int addr;
  /// 掩码数据/Mask Data
  unsigned char data[128];
} filter;

/**
 * @brief 自动模式下读取标签类型/The label type which read in automatic mode
 */
typedef struct tag_type_option
{
  /**
   * ISO18000-6B标签类型：0 - 不读取， 1 - 读取<br>
   *ISO18000-6B Tag Type: 0- do not read, 1 - read
   */
  int enable6c;
  /**
   * ISO18000-6C标签类型：0 - 不读取， 1 - 读取<br>
   * ISO18000-6C tag type: 0- do not read, 1 - read
   */
  int enable6b;
} tag_type_option;

/**
 * @brief 网络参数/network parameter
 */
typedef struct net_param
{
  /// IP
  char* ip;
  /// 子网掩码/subnet mask
  char* mask;
  /// 网关/gateway
  char* gw;
  /// 端口/port
  int port;
} net_param;

/**
 * @brief 自动上报格式字段标志/Automatic report format field flags
 *
 */
enum output_format_flags
{
  format_flag_epc = 1 << 0,
  format_flag_tid = 1 << 1,
  format_flag_user = 1 << 2,
  format_flag_rssi = 1 << 3,
  format_flag_ts = 1 << 4,
  format_flag_tag_type = 1 << 5,
  format_flag_ip = 1 << 6,
  format_flag_ant = 1 << 7,
  format_flag_custom4 = 1 << 26,
  format_flag_custom3 = 1 << 27,
  format_flag_custom2 = 1 << 28,
  format_flag_custom1 = 1 << 29,
  format_flag_custom0 = 1 << 30,
};

/**
 * @brief 自动模式上报格式选项/auto output data format settings
 */
typedef struct output_format
{
  /**
   * @brief 自动上报格式位图/format bitmap
   *
   * @details
   *
   * - 检查是否支持某个字段/checking flag
   *
   *@code
   * if(format.format & format_flag_ant) {
   *    printf("ant field enable");
   * }
   *@endcode
   *
   * - 设置某个字段/setting flag
   *
   * @code
   *   format.format = format.format | format_flag_ant | format_flag_epc |
   *                format_flag_ip | format_flag_tid;
   *
   *@endcode
   */
  int format;
  /// TID区起始地址/tid address
  int tid_addr;
  /// TID区长度/tid length
  int tid_len;
  /// USER区起始地址/ user address
  int user_addr;
  /// USER区长度/ user length
  int user_len;
} output_format;

/**
 * @brief 韦根设置参数/wiegand parameter
 */
typedef struct wg_option
{
  /// 韦根脉冲宽度。单位：10us / pulse width(10us)
  int width;
  /// 韦根脉冲间隔。单位：10us / pulse interval(10us
  int interval;
} wg_option;

/**
 * @brief 触发条件参数/trigger condition parameter
 */
typedef struct trigger_condition
{
  /// 触发电平/ trigger electrical level
  electrical_level level;
  /// 触发持续时间 / trigger duration
  int duration;
} trigger_condition;

/**
 * @brief 数据自动上报模式/auto output data report mode
 */
typedef enum report_type
{
  /**
   * 立即上报：读写器盘点一轮，上报一轮<br>
   * report right now: report after a round of inventory
   */
  report_now,
  /**
   * 定时上报（所有标签）：读写器在定时时间内进行盘点，定时时间到时将这段时间内盘点到的标签上报。<br>
   * report Periodically(all labels): inventory for a period and then report all labels read in this period
   */
  report_periodical1,
  /**
   * 定时上报（单标签）：针对单张标签的定时上报。例如上报的时间间隔设置为5秒，在第1秒的时候读到标签A，会立即上报该标签，但是往后的5秒时间内再次读到标签A时，不会再次上报。当时间超过5秒以后，再次识别到标签A，会立即上报。<br>
   * Timed reporting (single label): Timed reporting for a single label. For example, if the reporting interval is set
   * to 5 seconds, when tag A is read in the first second, the tag will be reported immediately, but if tag A is read
   * again in the next 5 seconds, it will not be reported again. When the time exceeds 5 seconds, tag A is recognized
   * again, and it will be reported immediately.
   */
  report_periodical2,
  /**
   *Modbus上报：此方式用于modbus协议 <br>
   * Modbus report: This method is used for modbus protocol.
   */
  report_modbus,
  /**
   * 实时上报：读写器在盘点过程中，识别到标签后实时上报标签的EPC或者UID。<br>
   * Real-time reporting: During the inventory process, the reader will report the EPC or UID of the tag in real time
   * after identifying the tag.
   */
  report_realtime,
} report_type;

/**
 * @brief 自动上报条件设置参数
 */
typedef struct report_condition
{
  /// 上报模式/reporting mode
  report_type rtype;
  /// 上报的时间间隔/reporting interval. Unit: seconds
  int interval;
} report_condition;

/**
 * @brief 心跳包设置参数
 */
typedef struct heartbeat_option
{
  /**
   * 是否开启心跳包。0--不开启，1--开启
   *
   * whether to open the heartbeat packet. 0--not open, 1--open.
   */
  int enable;
  /// 心跳包间隔/heartbeat packet interval.
  int interval;
  /**
   * 最长24个字节，心跳数据的ASCII码
   *
   *The maximum length of 24 bytes, the ASCII code of the heartbeat data
   */
  char* message;
} heartbeat_option;

/**
 * @brief 标签类型/rfid tag type
 */
typedef enum tag_type
{
  /// 6c
  tag6c,
  /// 6b
  tag6b
} tag_type;

typedef enum hf_tag_type
{
  hf_a,
  hf_b,
  hf_f,
  hf_v,
} hf_tag_type;

/**
 * @brief 标签的区域/memory bank of tag
 */
typedef enum memory_bank
{
  /// 密码区/password
  memory_bank_password,
  /// EPC
  memory_bank_epc,
  /// TID
  memory_bank_tid,
  /// USER
  memory_bank_user,
} memory_bank;

/**
 * @brief 一次标签读取的结果/The result of a reading
 *
 * 结果包含了epc号以及读取的数据（EPC，TID，USER中的某个区的数据，取决于读取时的参数）
 *
 * The result contains the EPC and the data that was read (EPC, TID, User, depending on the parameters at the
 * time of reading)
 */
typedef struct tag
{
  /// EPC
  unsigned char id[MAX_EPC_LEN + 1];
  /// EPC长度/EPC length
  int len;
  /// 读取的数据/data which read
  unsigned char data[MAX_USER_LEN + 1];
  /// 数据长度/data length
  int dataLen;
  /// 标签类型/tag type
  tag_type tagType;
} tag;

typedef struct hf_tag
{
  /// uid
  unsigned char uid[MAX_HF_UID_LEN + 1];
  /// uid length
  int uid_len;
  /// tag type
  hf_tag_type tagType;
} hf_tag;

typedef void (*tag_read_callback)(tag);

/**
 * @brief 锁定标签的区/Memory bank for locking
 */
typedef enum lock_memory_bank
{
  /// USER
  lock_memory_bank_user,
  /// TID
  lock_memory_bank_tid,
  /// EPC
  lock_memory_bank_epc,
  /// 访问密码区/access password memory bank
  lock_memory_bank_access_password,
  /// 销毁密码区/kill password memory bank
  lock_memory_bank_kill_password,
} lock_memory_bank;

/**
 * @brief 锁定类型/locking type
 */
typedef enum lock_type
{
  /// 开放/open
  lock_type_open,
  /// 锁定/lock
  lock_type_lock,
  /// 永久开放/open forever
  lock_type_open_forever,
  /// 永久锁定/lock forever
  lock_type_lock_forever,
} lock_type;

/**
 *@brief  wifi参数/wifi parameters
 */
typedef struct wifi_param
{
  /**
   *@brief wifi运行模式/wifi running mode
   *
   * 0--UDP Server，1--TCP Server
   */
  int mode;
  /// IP xxx.xxx.xxx.xxx
  char ip[MAX_IP_LEN + 1];
  /// 子网掩码/subnet mask
  char mask[MAX_IP_LEN + 1];
  /// 网关/gateway
  char gw[MAX_IP_LEN + 1];
  /// 通讯端口/communication port
  int port;
  /// wifi ssid
  char name[MAX_WIFI_SSID_LEN + 1];
  /// wifi密码/wifi password
  char password[MAX_WIFI_PASSWORD_LEN + 1];
  char remote_ip[MAX_IP_LEN + 1];
} wifi_param;

/**
 *@brief 超级网口配置信息/ TTL to Ethernet Modules parameters
 */
typedef struct super_net_if_param
{
  /**
   *@brief 运行模式/running mode
   *
   * 0--UDP Server，1--TCP Server
   */
  int mode;
  /// IP
  char* ip;
  /// 子网掩码/subnet mask
  char* mask;
  /// 网关/gateway
  char* gw;
  /// 通讯端口/communication port
  int port;
  char* remote_ip;
} super_net_if_param;

/**
 * @brief gsm网络参数/gsm4/5g parameter
 */
typedef struct four_g_param
{
  /**
   *@brief 运行模式/running mode
   *
   * 0--UDP Client，1--TCP Client
   */
  int mode;
  /// 通讯端口/communication port
  int server_port;
  /// IP
  char* server_ip;
} four_g_param;

#ifdef __cplusplus
extern "C"
{
#endif
  /**
   * @brief 暂定当前线程ms毫秒/pauses the current thread for ms Milliseconds
   * @param ms
   */
  __declspec(dllexport) void sleep_ms(int ms);

  /**
   * @brief 将字节数组转为16进制字符串/Converts a byte array to a hexadecimal string
   *
   * @param data 字节数组/byte array
   * @param len 字节数组长度/byte array length
   * @param [out] hex  输入/输出缓冲/input-output buffer
   */
  __declspec(dllexport) void to_hex_str(unsigned char* data, int len, char* hex);

  /**
   * @brief connect to reader via tcp/udp
   * @param ip ip
   * @param port port
   * @param protocol udp/tcp
   * @param [out] handle  connection handle
   * @return
   */
  __declspec(dllexport) int connect_net(char* ip, int port, transport_protocol protocol, int* handle);

  __declspec(dllexport) int connect_net_p32(char* ip, int port, transport_protocol protocol, int* handle);

  /**
   * @brief connect to reader via rs232/rs485
   * @param bd
   * @param address rs485 address
   * @param [out] handle
   * @return
   */
  __declspec(dllexport) int connect_rs(char*, enum baud bd, int address, int* handle);

  __declspec(dllexport) int connect_rs_p32(char*, enum baud bd, int address, int* handle);

  /**
   * @brief disconnect from reader
   * @param handle
   */
  __declspec(dllexport) void disconnect(int handle);

  /**
   * @brief get reader version information
   * @param handle
   * @param [out] version  version info
   *
   * @warning you should free the out version info memory: c_free(version)
   *
   * @code
   * char version[128];
   * get_version(handle, version);
   *
   * @endcode
   *
   * @return
   */
  __declspec(dllexport) int get_version(int handle, char* version);

  /**
   * @brief get baud rate of the port(rs232 or rs485)
   * @param handle
   * @param p port rs232/rs485
   * @param [out] b  baud rate
   * @return
   */
  __declspec(dllexport) int get_baud(int handle, port p, baud* b);

  /**
   * @brief set baud rate of the port(rs232 or rs485)
   * @param handle
   * @param p port rs232/rs485
   * @param b baud rate
   * @return
   */
  __declspec(dllexport) int set_baud(int handle, port p, baud b);

  /**
   * @brief get rs485 address
   * @param handle
   * @param [out] address  485 address
   * @return
   */
  __declspec(dllexport) int get_485_address(int handle, int* address);

  /**
   *@brief set rs485 address
   * @param handle
   * @param address 485 address
   * @return
   */
  __declspec(dllexport) int set_485_address(int handle, int address);

  /**
   * @brief change status of the relay
   * @param handle
   * @param relay relay number
   * @param status status of relay
   * @return
   */
  __declspec(dllexport) int set_relay_status(int handle, int relay, relay_status status);

  /**
   * @brief get status of the relay
   * @param handle
   * @param relay relay number
   * @param [out] status  status of relay
   * @return
   */
  __declspec(dllexport) int get_relay_status(int handle, int relay, relay_status* status);

  /**
   * @brief get relay count  of the reader
   * @param handle
   * @param [out] number  relay count
   * @return
   */
  __declspec(dllexport) int get_relay_number(int handle, int* number);

  /**
   *@brief 获取继电器自动控制参数/get relay control parameter
   * @param handle
   * @param relay 继电器编号/relay number
   * @param [out] option  控制参数/controls parameter
   * @return
   */
  __declspec(dllexport) int get_relay_option(int handle, int relay, relay_option* option);

  /**
   * @brief 设置继电器自动控制参数/set relay control parameter
   * @param handle
   * @param relay 继电器编号/relay number
   * @param option 控制参数/controls parameter
   * @return
   */
  __declspec(dllexport) int set_relay_option(int handle, int relay, relay_option option);

  /**
   * @brief 获取蜂鸣器状态/Gets the buzzer status
   * @param handle 连接句柄
   * @param [out] status   蜂鸣器状态: 0 开启 1 关闭 /  Buzzer status: 0 on 1 off
   * @return
   */
  __declspec(dllexport) int get_buzz(int handle, int* status);

  /**
   * @brief 设置蜂鸣器状态/sets the buzzer status
   * @param handle 连接句柄
   * @param status [in] 蜂鸣器状态: 0 开启 1 关闭/  Buzzer status: 0 on 1 off
   * @return
   */
  __declspec(dllexport) int set_buzz(int handle, int status);

  /**
   * @brief 设置工作模式/set working mode
   * @param handle 连接句柄
   * @param [in] mode 工作模式/working mode
   * @return 错误码
   */
  __declspec(dllexport) int set_work_mode(int handle, work_mode mode);

  /**
   * @brief 获取当前工作模式/get working mode
   * @param handle 连接句柄
   * @param [out] mode   工作模式/working mode
   * @return 错误码
   */
  __declspec(dllexport) int get_work_mode(int handle, work_mode* mode);

  /**
   * @brief 设置天线功率/Set antenna power
   * @param handle 连接句柄
   * @param  [in] ant 天线号/Antenna  NO
   * @param [in] power 功率/power
   * @return 错误码
   */
  __declspec(dllexport) int set_power(int handle, int ant, int power);

  /**
   * @brief 获取天线功率/Get antenna power
   * @param handle 连接句柄
   * @param ant 天线号/Antenna  NO
   * @param [out] power   功率/power
   * @return 错误码
   */
  __declspec(dllexport) int get_power(int handle, int ant, int* power);

  /**
   * @brief 设置模块频率区域/set rfid module frequency region
   * @param handle 连接句柄
   * @param  [out] region   频率区域/region
   * @return 错误码
   */
  __declspec(dllexport) int set_freq_region(int handle, freq_region region);

  /**
   * @brief 获取模块频率区域/get rfid module frequency region
   * @param handle 连接句柄
   * @param [in] region  频率区域/region
   * @return 错误码
   */
  __declspec(dllexport) int get_freq_region(int handle, freq_region* region);

  /**
   * @brief 自定义rfid模块频率参数/Customize RFID module frequency parameters
   * @param handle
   * @param option 频率参数/frequency parameter
   * @return
   */
  __declspec(dllexport) int set_user_freq(int handle, frequency_option option);

  /**
   * 获取rfid模块频率参数/get  frequency parameters of RFID module
   * @param handle
   * @param [out] option  频率参数/frequency parameter
   * @return
   */
  __declspec(dllexport) int get_user_freq(int handle, frequency_option* option);

  /**
   * @brief 获取自动模式下轮询的天线/Gets the polling antennas in automatic mode
   * @param handle 连接句柄
   * @param [out] ants   天线号数组/antenna NO.s
   * @param [out] count   天线个数/The number of antennas
   * @return 错误码
   */
  __declspec(dllexport) int get_inventory_ants(int handle, int* ants, int* count);

  /**
   * @brief 设置自动模式下轮询的天线/sets the polling antennas in automatic mode
   * @param handle 连接句柄
   * @param [in] ants 天线数组/antenna NO.s
   * @param [in] count 天线个数/The number of antennas
   * @return 错误码
   */
  __declspec(dllexport) int set_inventory_ants(int handle, const int* ants, int count);

  /**
   * @brief 恢复出厂设置/restore factory setting
   * @param handle 连接句柄
   * @return 错误码
   */
  __declspec(dllexport) int factory(int handle);

  /**
   * @brief 重启读写器/Restart
   * @param handle 连接句柄
   * @return 错误码
   */
  __declspec(dllexport) int reboot(int handle);

  /**
   * @brief 恢复wifi出厂设置/Restore WiFi factory Settings
   * @param handle 连接句柄
   * @return 错误码
   */
  __declspec(dllexport) int reset_wifi(int handle);

  /**
   * @brief 设置读写器时间/Set the reader time
   * @param handle 连接句柄
   * @param year 实际年-2000， 例如要设置2021年，需要传入2021-2000 = 21<br>
   * The actual year -2000, for example to set the year 2021, needs to pass in 2021-2000 = 21
   * @param month 月/month
   * @param day 日/day
   * @param hour 小时/hour
   * @param minute 分钟/minute
   * @param second 秒/second
   * @return 错误码
   */
  __declspec(dllexport) int set_time(int handle, int year, int month, int day, int hour, int minute, int second);

  /**
   * @brief 获取读写器时间/Gets the reader time
   *
   * @param [in]year 实际年等于该值+2000<br>The actual year is equal to this value plus 2000
   * @param [in]month 月/month
   * @param [in]day 日/day
   * @param [in]hour 小时/hour
   * @param [in]minute 分钟/minute
   * @param [in]second 秒/second
   * @return 错误码
   */
  __declspec(dllexport) int get_time(int handle, int* year, int* month, int* day, int* hour, int* minute, int* second);

  /**
   * @brief 设置rfid标签算法相关参数/Set the RFID tag algorithm parameters
   * @param handle
   * @param algorithm 参数/parameter
   * @return
   */
  __declspec(dllexport) int set_tag_algorithm(int handle, tag_algorithm algorithm);

  /**
   * @brief 获取rfid标签算法相关参数/get the RFID tag algorithm parameters
   * @param handle
   * @param [out] algorithm  参数/parameter
   * @return
   */
  __declspec(dllexport) int get_tag_algorithm(int handle, tag_algorithm* algorithm);

  /**
   * @brief 设置标签过滤器/Set Label Filters
   *
   * @param handle
   * @param flt 过滤器参数/filter parameter
   * @return
   */
  __declspec(dllexport) int set_filter(int handle, filter flt);

  /**
   * @brief 获取标签过滤器/get Label Filters
   *
   * @param handle
   * @param [out] flt  过滤器参数/filter parameter
   * @return
   */
  __declspec(dllexport) int get_filter(int handle, filter* flt);

  /**
   * @brief 设置自动模式下读取标签类型/Sets the read label type in automatic mode
   * @param handle
   * @param option
   * @return
   */
  __declspec(dllexport) int set_tag_type_option(int handle, tag_type_option option);

  /**
   * @brief 获取自动模式下读取标签类型/gets the read label type in automatic mode
   * @param handle
   * @param [out] option
   * @return
   */
  __declspec(dllexport) int get_tag_type_option(int handle, tag_type_option* option);

  /**
   * @brief 获取网络参数/Get network parameters
   * @param handle
   * @param [out] param  网络参数/network parameter
   *
   * @warning 使用完param之后，请释放内存free_network_param(&param);<br>
   * After using param, free the memory with free_network_param(&amp;param)
   * @return
   */
  __declspec(dllexport) int get_network_param(int handle, net_param* param);

  /**
   * @brief 释放网络参数结构内存/Free network parameter structure memory
   * @param param
   */
  __declspec(dllexport) void free_network_param(net_param* param);

  /**
   * @brief 设置网络参数/set network parameters
   * @param handle
   * @param param 网络参数/network parameter
   * @return
   */
  __declspec(dllexport) int set_network_param(int handle, net_param param);

  /**
  * @brief 获取mac地址/Get MAC address
  *
  * 基本用法
  * @code
  char *mac = "12:34:89:09:aa:bb";
  int ret = set_mac(id, mac);
  assert(ret == err_ok);

  mac = NULL;
  ret = get_mac(id, &mac);
  assert(ret == err_ok);
  assert(strcmp(mac, "12:34:89:09:aa:bb") == 0);

  c_free(mac);
  * @endcode
  *
  * 注意释放内存
  * @param handle
  * @param [out] mac
  * @return
  */
  __declspec(dllexport) int get_mac(int handle, char** mac);

  /**
   * 释放内存, free函数的包装<br>
   * Free memory, free function wrapper
   * @param pointer
   */
  __declspec(dllexport) void c_free(void* pointer);

  /**
   * @brief 设置mac地址/set mac address
   * @param handle
   * @param mac mac
   * @return
   */
  __declspec(dllexport) int set_mac(int handle, char* mac);

  /**
   * @brief 获取上报接口状态/get report port status
   * @param handle 连接句柄
   * @param p 接口/port
   * @param [out] status  0 - 未启用  1 - 已启用
   *
   * 0 - Not Enabled 1 - Enabled
   *
   * @return 错误码
   */
  __declspec(dllexport) int get_report_port_status(int handle, port p, int* status);

  /**
   * @brief 设置上报接口状态/set report port status
   * @param handle 连接句柄
   * @param p 接口/port
   * @param status  0 表示未启用  1 表示已启用
   *
   * 0 - Not Enabled 1 - Enabled
   * @return 错误码
   */
  __declspec(dllexport) int set_report_port_status(int handle, port p, int status);

  /**
   * @brief 设置自动输出格式/set auto output data format
   * @param handle
   * @param format 格式参数/format parameter
   * @return
   */
  __declspec(dllexport) int set_output_format(int handle, output_format format);

  /**
   * @brief 获取自动输出格式/get auto output data format
   * @param handle
   * @param [out] format  格式参数/format parameter
   * @return
   */
  __declspec(dllexport) int get_output_format(int handle, output_format* format);

  /**
   * @brief 设置自定义上报字段内容/set custom report content
   * @param handle
   * @param seq 字段序号/field number
   * @param content 内容/content
   * @return
   */
  __declspec(dllexport) int set_custom_report_content(int handle, int seq, char* content);

  /**
   * @brief 获取自定义上报内容/get custom report content
   * @param handle 连接句柄
   * @param seq 序号/field number
   * @param [out] content  上报内容/content
   * @return 错误码
   */
  __declspec(dllexport) int get_custom_report_content(int handle, int seq, char** content);

  /**
   * @brief 设置韦根参数/set wiegand parameter
   * @param handle
   * @param option 参数/parameter
   * @return
   */
  __declspec(dllexport) int set_wg_option(int handle, wg_option option);

  /**
   * @brief 获取韦根参数/get wiegand parameter
   * @param handle
   * @param [out] option  参数/parameter
   * @return
   */
  __declspec(dllexport) int get_wg_option(int handle, wg_option* option);

  /**
   *@brief 获取触发条件参数/get trigger condition param
   * @param handle
   * @param seq 触发器序号/trigger number
   * @param condition 参数/param
   * @return
   */
  __declspec(dllexport) int get_trigger_option(int handle, int seq, trigger_condition* condition);

  /**
   *@brief 设置触发条件参数/set trigger condition param
   * @param handle
   * @param seq 触发器序号/trigger number
   * @param [out] condition  参数/param
   * @return
   */
  __declspec(dllexport) int set_trigger_option(int handle, int seq, trigger_condition condition);

  /**
   * @brief 获取读写器ID/get reader id
   * @param handle
   * @param [out] id  id
   *
   * @warning free id after you will never use it any more c_free(id)
   * @return
   */
  __declspec(dllexport) int get_id(int handle, char** id);

  /**
   * @brief 获取读写器名字/get reader name
   * @param handle
   * @param [out] mode  名字/name
   *
   * @warning free mode after you will never use it any more c_free(mode)
   *
   * @return
   */
  __declspec(dllexport) int get_reader_name(int handle, char** mode);

  /**
   *@brief 设置读写器名字/set reader name
   * @param handle
   * @param mode 名字/name
   * @return
   */
  __declspec(dllexport) int set_reader_name(int handle, char* mode);

  /**
   * @brief 设置自动上报条件设置参数/set auto output report condition param
   * @param handle
   * @param condition 条件参数/condition param
   * @return
   */
  __declspec(dllexport) int set_report_condition(int handle, report_condition condition);

  /**
   * @brief 获取自动上报条件设置参数/get auto output report condition param
   * @param handle
   * @param [out] condition  条件参数/condition param
   * @return
   */
  __declspec(dllexport) int get_report_condition(int handle, report_condition* condition);

  /**
   * @brief 设置心跳参数/set heartbeat parameter
   * @param handle
   * @param option 心跳参数/heartbeat parameter
   * @return
   */
  __declspec(dllexport) int set_heartbeat_option(int handle, heartbeat_option option);

  /**
   * @brief 获取心跳参数/get heartbeat parameter
   * @param handle
   * @param [out] option  心跳参数/heartbeat parameter
   *
   *  @warning free option after you will never use it any more with free_heartbeat_option(option)
   *
   * @return
   */
  __declspec(dllexport) int get_heartbeat_option(int handle, heartbeat_option* option);

  /**
   * @brief 释放心跳参数内存/release  heartbeat parameter memory
   * @param option
   */
  __declspec(dllexport) void free_heartbeat_option(heartbeat_option* option);

  /**
   * @brief 设置当前工作天线/Set the current working antenna
   * @param handle
   * @param ant 天线序号/The antenna number
   * @return
   */
  __declspec(dllexport) int set_work_ant(int handle, int ant);

  /**
   * @brief 获取读写器天线个数/Get the number of reader antennas
   * @param handle
   * @param [out] count  天线个数/The number of antennas
   * @return
   */
  __declspec(dllexport) int get_ant_count(int handle, int* count);

  /**
   * @brief 列出附近的标签/List nearby tags
   *
   * 基本用法参考首页.<p>
   * Refer to the home page for basic usage
   *
   * @note 当memory参数是epc时， addr,len参数无效，会列出epc<p>
   * When the memory parameter is EPC, the addr,len parameters are invalid and the EPC is listed
   *
   * @param handle
   * @param memory 指定需要列出的内存区/memory that need to be listed
   * @param addr 起始地址/start address
   * @param len 读取长度/length
   * @param password 访问密码/access password
   * @param passwordLen 密码长度/password length
   * @param [out] tags  tags
   * @param [out] count  标签个数/tag count
   *
   *  @warning free tags after you will never use it any more free_tags(tags)
   * @return
   */
  __declspec(dllexport) int list6c(
      int handle,
      memory_bank memory,
      int addr,
      int len,
      unsigned char* password,
      int passwordLen,
      tag* tags,
      int* count);

  __declspec(dllexport) int list6c_with_callback(
      int handle,
      memory_bank memory,
      int addr,
      int len,
      unsigned char* password,
      int passwordLen,
      tag_read_callback callback);

  /**
   *@brief 读取标签的特定数据/Read the specific data of the label
   *
   * 参考主页用法
   *
   * @param handle
   * @param memory 指定需要读取的内存区/memory that need to be read
   * @param addr 起始地址/start address
   * @param len 长度/length
   * @param epc 待读取的epc/The EPC to be read
   * @param epcLen epc长度/length
   * @param password 访问密码/access password
   * @param passwordLen 密码长度/password length
   * @param [out] data  读到的数据/data read
   * @param [out] count  数据长度/data length
   * @return
   */
  __declspec(dllexport) int read6c(
      int handle,
      memory_bank memory,
      int addr,
      int len,
      unsigned char* epc,
      int epcLen,
      unsigned char* password,
      int passwordLen,
      unsigned char* data,
      int* count);

  /**
   *@brief 往标签的特定区域写入数据/Writes data to a specific area of the label
   *
   *
   * @param handle
   * @param memory 指定需要写入的内存区/memory that need to be writen
   * @param addr 起始地址/start address
   * @param epc 待写入标签的epc/The EPC to be writen
   * @param epcLen epc长度/length
   * @param data 待写入数据/data to be writen
   * @param dataLen 数据长度/data length
   * @param password 访问密码/access password
   * @param passwordLen 密码长度/password length
   * @return
   */
  __declspec(dllexport) int write6c(
      int handle,
      memory_bank memory,
      int addr,
      unsigned char* epc,
      int epcLen,
      unsigned char* data,
      int dataLen,
      unsigned char* password,
      int passwordLen);

  /**
   * @brief 锁定标签/locking tag
   * @param handle
   * @param lockMemory 锁定区域/locking memory
   * @param lt 锁定类型/locking type
   * @param epc 待锁定标签epc/tag to be locked
   * @param epcLen epc长度/epc length
   * @param password 访问密码/access password
   * @param passwordLen 密码长度/password length
   * @return
   */
  __declspec(dllexport) int lock6c(
      int handle,
      lock_memory_bank lockMemory,
      lock_type lt,
      unsigned char* epc,
      int epcLen,
      unsigned char* password,
      int passwordLen);

  /**
   *@brief 毁灭标签/killing tag
   * @param handle
   * @param epc 待毁灭标签/tag to be killed
   * @param epcLen 标签长度/epc length
   * @param killPassword 毁灭密码/kill password
   * @param passwordLen 密码长度/password length
   * @return
   */
  __declspec(
      dllexport) int kill6c(int handle, unsigned char* epc, int epcLen, unsigned char* killPassword, int passwordLen);

  /**
   * @brief 快速写epc/fast write epc
   * @param handle
   * @param oldEpc 旧的epc，如果这个参数不为空，则更新这个标签的epc，否则随机写一个epc。<p>
   * if oldEpc is null, it will write a random epc. or it will update the old epc to new epc.
   * @param oldEpcLen 旧epc长度/old epc length
   * @param newEpc 待写入epc/epc
   * @param newEpcLen epc长度/epc length
   * @param password 访问密码/access password
   * @param passwordLen 密码长度/password length
   * @return
   */
  __declspec(dllexport) int quick_write_epc(
      int handle,
      unsigned char* oldEpc,
      int oldEpcLen,
      unsigned char* newEpc,
      int newEpcLen,
      unsigned char* password,
      int passwordLen);

  /**
   * @brief 列出附近的6b标签/List nearby 6B tags
   * @param handle
   * @param [out] tags  tags
   * @param [out] count
   *
   *  @warning free tags after you will never use it any more free_tags(tags)
   * @return
   */
  __declspec(dllexport) int list6b(int handle, tag* tags, int* count);

  /**
   * @brief 读取6b标签特定数据/read 6b data
   * @param handle
   * @param uid uid
   * @param uidLen uid长度/uid length
   * @param addr 起始地址/start address
   * @param len 待读取长度/data length to read
   * @param [out] data  读取数据/data
   * @param [out] count  数据长度/data length
   * @return
   */
  __declspec(dllexport) int read6b(
      int handle,
      unsigned char* uid,
      int uidLen,
      int addr,
      int len,
      unsigned char* data,
      int* count);

  /**
   * @brief 写入6b标签数据/write 6b data
   *
   * @param handle
   * @param uid uid
   * @param uidLen uid长度/uid length
   * @param addr 起始地址/start address
   * @param data 待写入数据/data to be write
   * @param dataLen 数据长度/data length
   * @return
   */
  __declspec(
      dllexport) int write6b(int handle, unsigned char* uid, int uidLen, int addr, unsigned char* data, int dataLen);

  /**
   * @brief 锁定6b标签/lock 6b tag
   * @param handle
   * @param uid uid
   * @param uidLen  uid长度/uid length
   * @param addr  起始地址/start length
   * @return
   */
  __declspec(dllexport) int lock6b(int handle, unsigned* uid, int uidLen, int addr);

  /**
   *@brief 设置标签报警器/set tag alarm
   * @param handle
   * @param flt 报警器参数/alarm parameter
   * @return
   */
  __declspec(dllexport) int set_alarm(int handle, filter flt);

  /**
   *@brief 获取标签报警器/get tag alarm
   * @param handle
   * @param [out] flt 报警器参数/alarm parameter
   * @return
   */
  __declspec(dllexport) int get_alarm(int handle, filter* flt);

  /**
   *@brief 设置快速读取TID功能/set fast tid
   * @param handle
   * @param enable 1 enable 0 disable
   * @return
   */
  __declspec(dllexport) int set_fast_tid(int handle, int enable);

  /**
   *@brief 获取快速读取TID功能/get fast tid
   * @param handle
   * @param [out] enable 1 enable 0 disable
   * @return
   */
  __declspec(dllexport) int get_fast_tid(int handle, int* enable);

  /**
   * @brief 获取触发器个数/Gets the number of triggers
   * @param handle
   * @param [out] number  触发器个数/number of trigger
   * @return
   */
  __declspec(dllexport) int get_trigger_number(int handle, int* number);

  /**
   * @brief 监听自动模式下读到的数据/Listen for data read in automatic mode
   *
   *
   * @param handle
   * @param callback
   * @return
   */
  __declspec(dllexport) int start_listen_auto_read(int handle, auto_read_callback callback);

  __declspec(dllexport) int start_listen_auto_read_sync(int handle, auto_read_callback callback);

  /**
   * @brief 停止监听数据/Stop listening for data
   * @param handle
   * @return
   */
  __declspec(dllexport) int stop_listen_auto_read(int handle);

  /**
   * @brief 监听自动模式下udp上报的数据/Listen to data reported by UDP in automatic mode.
   *
   * 部分型号读写器udp连接下需要用此函数监听
   *
   *  Some types of readers need to listen to this function under UDP connection
   *
   * @param ip 监听ip地址
   * @param port 监听端口
   * @param callback 有输出的时候会回调该函数
   * @return
   */
  __declspec(dllexport) int start_listen_udp_auto_read(char* ip, int port, auto_read_callback callback);

  __declspec(dllexport) int start_listen_udp_auto_read_p32(int port, auto_read_callback callback);

  /**
   * @brief 停止监听udp数据/stop listen auto output data
   * @param ip 监听ip地址/reader ip
   * @param port 监听端口/reader port
   * @return
   */
  __declspec(dllexport) int stop_listen_udp_auto_read(char* ip, int port);

  __declspec(dllexport) int stop_listen_udp_auto_read_p32(int port);

  /**
   * @brief start to scan device
   *
   * scan device(must be in the same network) which is online, for example
   *
   * @code
   * void example_callback(rfid_device device)
   * {
   *      printf("device ip[%s] port[%d]\n", device.ip, device.port);
   *      free_rfid_device(device);
   * }
   *
   * you should free data by yourself if you don't need the data anymore.
   *
   * you should stop scanning also by yourself with @ref stop_device_scan
   *
   * @endcode
   *
   * @param callback
   * @return
   */
  __declspec(dllexport) int start_device_scan(rfid_device_scan_callback callback);

  /**
   * @brief free memory from start_device_scan api
   *
   * @see start_device_scan
   *
   * @param device
   */
  __declspec(dllexport) void free_rfid_device(rfid_device device);

  /**
   * stop device scan request
   * @return
   */
  __declspec(dllexport) int stop_device_scan();

  /**
   * @brief 获取设备型号/get equipment mode
   * @param handle
   * @param [out] mode 型号编号/equipment mode
   * @return
   */
  __declspec(dllexport) int get_equipment_mode(int handle, int* mode);

  /**
   * @brief 查询读写器是否装配了wifi模块/get wifi module availability
   * @param handle
   * @param [out] available
   * @return
   */
  __declspec(dllexport) int is_wifi_available(int handle, int* available);

  /**
   * @brief get_modul_type
   * @param handle
   * @param [out] modul_type
   * @return
   */
  __declspec(dllexport) int get_modul_type(int handle, int* modul_type);

  /**
   * @brief 查询读写器是否装配了gsm4/5g模块/get gsm4g/5g module availability
   * @param handle
   * @param [out] available
   * @return
   */
  __declspec(dllexport) int is_4g_available(int handle, int* available);

  /**
   * @brief 查询读写器是否装配了韦根/get wiegand module availability
   * @param handle
   * @param [out] available
   * @return
   */
  __declspec(dllexport) int is_wg_available(int handle, int* available);

  /**
   * @brief 查询读写器是否装配了超级网口/get super network interface module availability
   * @param handle
   * @param [out] available
   * @return
   */
  __declspec(dllexport) int is_super_net_if_available(int handle, int* available);

  /**
   * @brief 设置wifi参数/set wifi parameter
   * @param handle
   * @param param  wifi参数/wifi parameter
   * @return
   */
  __declspec(dllexport) int set_wifi_param(int handle, wifi_param param);

  /**
   * @brief 获取wifi参数/get wifi parameter
   * @param handle
   * @param [out] param WiFi参数/wifi parameter
   *
   * @warning 使用完毕之后请使用free_wifi_param释放param的内存
   * @return
   */
  __declspec(dllexport) int get_wifi_param(int handle, wifi_param* param);

  /**
   * @brief 设置超级网口参数/set super network interface parameter
   * @param handle
   * @param param  超级网口参数/super network interface parameter
   * @return
   */
  __declspec(dllexport) int set_super_net_if_param(int handle, super_net_if_param param);

  /**
   * @brief 获取超级网口参数/get super network interface parameter
   * @param handle
   * @param [out] param 超级网口参数/super network interface parameter
   *
   * @warning 使用完毕之后请使用free_super_net_if_param释放param的内存
   * @return
   */
  __declspec(dllexport) int get_super_net_if_param(int handle, super_net_if_param* param);

  /**
   * @brief 释放内存/free memory
   *
   * @param param
   */
  __declspec(dllexport) void free_super_net_if_param(super_net_if_param* param);

  /**
   * @brief 设置gsm参数/set gsm parameter
   * @param handle
   * @param param  gsm参数/gsm parameter
   * @return
   */
  __declspec(dllexport) int set_4g_param(int handle, four_g_param param);
  /**
   * @brief 获取gsm参数/get gsm parameter
   * @param handle
   * @param [out] param gsm参数/gsm parameter
   *
   * @warning 使用完毕之后请使用free_4g_param释放param的内存
   * @return
   */
  __declspec(dllexport) int get_4g_param(int handle, four_g_param* param);

  /**
   * @brief 释放内存
   *
   * @param param
   */
  __declspec(dllexport) void free_4g_param(four_g_param* param);

  /**
   * @brief 获取核心库版本
   *
   * @return  版本信息
   *
   * @warning 注意使用完返回值之后释放内存
   */
  __declspec(dllexport) char* core_version();

  /**
   * @brief 获取错误描述/get error message description
   *
   * @param code 错误码/error code
   * @return 错误描述信息/error description
   *
   * @warning 注意使用完返回值之后释放内存/free result memory after use
   */
  __declspec(dllexport) void error_description(int code, char* errMessage);

  /**
   * @brief 进入/退出调试模式 enable/disable debug mode
   *
   * @param enable
   */
  __declspec(dllexport) void set_debug(int enable, int toFile);

  __declspec(dllexport) int start_s108_proxy();

  __declspec(dllexport) void stop_s108_proxy();

  /**
   * @brief connect to r2000 reader via tcp
   * @param ip ip
   * @param port port
   * @param protocol udp/tcp
   * @param [out] handle  connection handle
   * @return
   */
  __declspec(dllexport) int connect_r2000_tcp(char* ip, int port, int device_type, int* handle);

  /**
   * @brief connect to r2000 reader via udp
   * @param ip
   * @param port
   * @param handle
   * @return
   */
  __declspec(dllexport) int connect_r2000_udp(char* ip, int port, int device_type, int* handle);

  /**
   * @brief connect to reader via rs232/rs485
   * @param port serialport
   * @param bd
   * @param address rs485 address, use 0xff for common use
   * @param [out] handle
   * @return
   */
  __declspec(dllexport) int connect_r2000_rs(char* port, enum baud bd, int address, int device_type, int* handle);

  __declspec(dllexport) int connect_p218_tcp(char* ip, int port, int device_type, int* handle);
  __declspec(dllexport) int connect_p218_udp(char* ip, int port, int device_type, int* handle);

  /**
   * @brief connect_p218_rs
   * @param port
   * @param bd
   * @param address 地址，如果是<=0,表示不带地址参数
   * @param device_type
   * @param handle
   * @return
   */
  __declspec(dllexport) int connect_p218_rs(char* port, enum baud bd, int address, int device_type, int* handle);
  __declspec(dllexport) int connect_p218_usb(int device_type, int* handle);

  __declspec(dllexport) int get_report_option(int handle, report_option* option_udp, report_option* option_tcp);

  __declspec(dllexport) int set_report_option(int handle, report_option option_udp, report_option option_tcp);

  // internal api
  __declspec(dllexport) char* internal_get_ip_from_udp_report_option(report_option* option);

  __declspec(dllexport) void set_log_level(int level);
  __declspec(dllexport) void set_log_output_file_path(char* path);

  /**
   * @brief list hf tags
   * @param handle connection handle
   * @param [out] tags
   * @param [out] count number of tags read
   * @return error number
   */
  __declspec(dllexport) int hf_list_tag(int handle, hf_tag* tags, int* count);

  /**
   * @brief read 15693 block data
   * @param handle connection handle
   * @param uid uid of the tag, if not specified, it will choose a random 15693 tag to read
   * @param uid_len uid length
   * @param address where to start read
   * @param block_count number of block to be read
   * @param [out] data data that read
   * @param [out] data_len data length that read
   * @return error code
   */
  __declspec(dllexport) int hf_read_15693_block(
      int handle,
      unsigned char* uid,
      int uid_len,
      int address,
      int block_count,
      unsigned char* data,
      int* data_len);

  /**
   * @brief write block data
   * @param handle connection handle
   * @param uid uid of the tag, if not specified, it will choose a random 15693 tag to read
   * @param uid_len uid length
   * @param address where to start write
   * @param data data that read(only 4 bytes allowed)
   * @param data_len data length that read(fixed 4)
   * @return error code
   */
  __declspec(dllexport) int hf_write_15693_block(
      int handle,
      unsigned char* uid,
      int uid_len,
      int address,
      unsigned char* data,
      int data_len);

  __declspec(dllexport) int hf_write_15693_multi_block(
      int handle,
      unsigned char* uid,
      int uid_len,
      int address,
      unsigned char* data,
      int data_len);

  /**
   * @brief read 14443 block data
   * @param handle connection handle
   * @param uid uid of the tag, if not specified, it will choose a random 14443 tag to read
   * @param uid_len uid length
   * @param address where to start read
   * @param block_count number of block to be read
   * @param [out] data data that read
   * @param [out] data_len data length that read
   * @return error code
   */
  __declspec(dllexport) int hf_read_14443_block(
      int handle,
      unsigned char* uid,
      int uid_len,
      int address,
      int block_count,
      unsigned char* data,
      int* data_len);

  /**
   * @brief write block data
   * @param handle connection handle
   * @param uid uid of the tag, if not specified, it will choose a random 14443 tag to read
   * @param uid_len uid length
   * @param address where to start write
   * @param data data that read(only 4 bytes allowed)
   * @param data_len data length that read(fixed 4)
   * @return error code
   */
  __declspec(dllexport) int hf_write_14443_block(
      int handle,
      unsigned char* uid,
      int uid_len,
      int address,
      unsigned char* data,
      int data_len);

  __declspec(dllexport) int hf_write_14443_multi_block(
      int handle,
      unsigned char* uid,
      int uid_len,
      int address,
      unsigned char* data,
      int data_len);

  /**
   * @brief select 14443A fan, default 0
   * @param handle connection handle
   * @param uid  uid of the tag, if not specified, it will choose a random 14443 tag to read
   * @param uid_len uid length
   * @param fan fan number
   * @return error code
   */
  __declspec(dllexport) int hf_select_14443_fan(int handle, unsigned char* uid, int uid_len, int fan);

#ifdef __ANDROID__
  int connect_local_s108(int* handle);
  void set_s108_wifi_led(int on);
  int get_s108_wifi_led();
  void set_s108_4g_led(int on);
  int get_s108_4g_led();
  void set_s108_com_led(int on);
  int get_s108_com_led();
#endif

  // ibat2000 connect
  __declspec(dllexport) int ibat2000_rs232_connect(char* port, baud bd, int* handle, int timeout);

  // p218 desktop reader keyboard simulator api
  __declspec(dllexport) int set_sk_output_prefix(int handler, int enable, char* prefix);
  __declspec(dllexport) int get_sk_output_prefix(int handler, int* enable, char* prefix);
  __declspec(dllexport) int set_sk_output_format(int handler, int format);
  __declspec(dllexport) int get_sk_output_format(int handler, int* format);
  __declspec(dllexport) int set_sk_output_addr(int handler, int addr, int len);
  __declspec(dllexport) int get_sk_output_addr(int handler, int* addr, int* len);
  __declspec(dllexport) int set_sk_output_interval(int handler, int interval);
  __declspec(dllexport) int get_sk_output_interval(int handler, int* interval);
  __declspec(dllexport) int set_sk_new_line_enter(int handler, int enable);
  __declspec(dllexport) int get_sk_new_line_enter(int handler, int* enable);
  __declspec(dllexport) int get_sk_output_memory_bank(int handler, int* bank);
  __declspec(dllexport) int set_sk_output_memory_bank(int handler, int bank);

  // eas api
  __declspec(dllexport) int change_eas_state(
      int handler,
      unsigned char* epc,
      int epc_len,
      unsigned char* password,
      int password_len,
      int status);

  __declspec(dllexport) int eas_alarm(int handler);

  // p218 start auto mode tmp
  __declspec(dllexport) int start_auto_mode_tmp(int handler);
  // p218 stop auto mode tmp
  __declspec(dllexport) int stop_auto_mode_tmp(int handler);

#ifdef __cplusplus
}
#endif

#endif
