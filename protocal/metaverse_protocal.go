using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Threading;
using Microsoft.Win32;
using System.Net;
using System.Net.Sockets;
using System.Collections;
using System.Collections.Generic;
using System.Diagnostics;
using System.Text;

namespace SerialCom
{
    class metaverse_protocal
    {

        public Byte[] Uart_Rx_buf = new Byte[256];
        public Int32 Uart_Rx_idx = 0;
        private Byte meta_api_state = 0;
        private Byte meta_api_cmd_type = 0;
        private Byte meta_api_extend_len = 0;
        private Byte meta_api_data_type = 0;
        private UInt32 meta_api_buf_start_idx = 0;
        private Byte[] meta_api_support_versions = new byte[] { 0,1,2,16};
        private Byte meta_api_min_support_version = 2;
        private Byte meta_api_max_support_version = 16;
        private Int16 meta_api_version = 0;
        private UInt16 meta_api_len = 0;
        private UInt16 meta_api_data_len = 0;
        private UInt16 meta_api_rx_idx = 0;
        private UInt16 meta_api_seq = 0;
        private UInt16 meta_api_cmdset = 0;
        private UInt16 meta_api_cmdid = 0;
        private const int meta_api_buf_v0_max_size = 63;
        private const int meta_api_buf_v1_max_size = 63;
        private const int meta_api_buf_v2_max_size = 4095;
        private const int meta_api_buf_v16_max_size = 65535;
        //if you only support version 2,just change blow.
        private const int meta_api_buf_max_size = meta_api_buf_v16_max_size +1;
        private Byte[] meta_api_data_buf = new Byte[meta_api_buf_max_size];
        private Byte[] meta_api_buf = new Byte[meta_api_buf_max_size];
        private int meta_api_buf_idx = 0;
        private UInt16 meta_api_extend_idx = 0;
        private List<byte> meta_api_package_buffer = new List<byte>();
        private List<byte> meta_api_data_buffer = new List<byte>();
        private int total_package_count = 0;
        private int total_packege_len2 = 0;

        private string last_save_data;
        private long last_time_ticks = 0;

        private UInt16 meta_api_crc16 = 0;
        private UInt16 meta_api_crc16_calc = 0;
        private UInt32 meta_api_crc32 = 0;
        private UInt32 meta_api_crc32_calc = 0;

        private crc16 crc16_obj = new crc16();
        private crc32 crc32_obj = new crc32();
        private bool Frm_main_init = false;
        public metaverse_unreal metaverse_unreal_obj = new metaverse_unreal();

        private long sendStartTime = DateTime.Now.ToUniversalTime().Ticks;

        public MainForm mainform=null;

        public void init(string local_ip1, string remote_ip1, int remote_port1)
        {
            metaverse_unreal_obj.send_init(local_ip1, remote_ip1, remote_port1);
        }

        public void send_metaverse_unreal_data(UInt32 parm_a, UInt16 parm_b,float parm_c)
        {

            var ts = DateTime.Now.ToUniversalTime().Ticks;

            if (ts - sendStartTime > 30 * 10000)
            {
                sendStartTime = ts;
                metaverse_unreal_obj.send_metaverse_unreal_data((float)parm_a / 1000, parm_b, parm_c);
            }

        }
        #region 协议执行
        private int Process_metaverse_packet(Byte version,UInt16 meta_api_cmdset, UInt16 meta_api_id, UInt16 SEQ, 
            Byte[] meta_api_packet_buf, UInt16 meta_api_packet_len,  Byte[] meta_api_data_buf, UInt16 meta_data_len)
        {
            Byte[] data_buf_all=new Byte[1];
            Byte[] package_buf_all = new Byte[1];
            int total_packege_len = meta_data_len;
            

            if (version < 2)
            {
                data_buf_all = meta_api_data_buf;
                total_packege_len2 = meta_data_len;
                total_package_count++;
                meta_api_package_buffer.AddRange(meta_api_packet_buf);
            }
            else
            {
                if ((SEQ & 0x8000) > 0)
                {
                    meta_api_data_buffer.AddRange(meta_api_data_buf);
                    meta_api_package_buffer.AddRange(meta_api_packet_buf);
                    data_buf_all = meta_api_data_buffer.ToArray();
                    meta_api_data_buffer.Clear();
                    total_packege_len2 += meta_data_len;
                    total_packege_len = data_buf_all.Length;
                    total_package_count++;
                    /*
                    if (version == 16 )
                    {
                        for (int i = 0; i < meta_data_len; i++)
                        {
                            if (data_buf_all[i] != (Byte)i)
                            {
                                Console.WriteLine("seems error");
                            }
                        }
                    }*/
                }
                else
                {
                    total_package_count++;
                    total_packege_len2 += meta_data_len;
                    meta_api_data_buffer.AddRange(meta_api_data_buf);
                    meta_api_package_buffer.AddRange(meta_api_packet_buf);
                    return 0;
                }
            }

            if (meta_api_cmdset == 0 )
            {
                send_metaverse_unreal_data(1, 2, 3);
                bool recieve_ok = (version < 2) || ((version >= 2) && (total_package_count - 1 == (SEQ & 0x7FFF)));
                var spritfstr = "";
                if ( recieve_ok )
                {
                   spritfstr = String.Format(@"recvok ! version  {0:D} ,received len {1:D} , packege_count {2:D}", version, total_packege_len2, total_package_count);
                }
                else
                {
                    spritfstr = String.Format(@"some package lost ! version  {0:D} ,received len {1:D} , packege_count {2:D}", version, total_packege_len2, total_package_count);
                }
                mainform.textBoxReceive.Text += spritfstr + "\r\n"; ;
                //mainform.textBoxReceive.SelectionStart = mainform.textBoxReceive.Text.Length;
                //mainform.textBoxReceive.ScrollToCaret();//滚动到光标处
            }

            save_packet(meta_api_package_buffer.ToArray());
            meta_api_data_buffer = new List<byte>();
            meta_api_package_buffer = new List<byte>();
            total_packege_len2 = 0;
            total_package_count = 0;
            return 0;
        }

        #endregion
        /*
        public class metaverseProtocalS
        {
            Byte version;
            Byte DataType;
            Byte CmdType; 
            Byte ENC;
            UInt16 meta_api_cmdset;
            UInt16 meta_api_id;
            Byte[] meta_api_data_buf;
            UInt16 meta_data_len;
        };*/
        public Byte[] metaverseProtocalGenMassive(Byte version, Byte data_type, Byte cmd_type, Byte ENC, UInt16 cmd_set, UInt16 cmd_id,  Byte[] data_buf, UInt32 data_len)
        {
            Byte[] api_buf;
            UInt16 max_payload = meta_api_buf_v16_max_size - 22;
            //not recommend version <=2 to gen massive buf
            if (version == 2)
            {
                max_payload = meta_api_buf_v2_max_size - 13;
            }
            Byte[] api_buf_all = new byte[(data_len/ max_payload+1) *65535];
            int api_buf_all_idx = 0;
            List<byte> data_buffer = new List<byte>(data_buf.Length);
            data_buffer.AddRange(data_buf);
            //Byte[] api_buf_all ;
            List<byte> buffer = new List<byte>();
            
            UInt16 SEQ = 0x0000;
            int i = 0;
            int payload_len_all = 0;
            for (;i< data_len; i+= max_payload)
            {
                UInt16 payload_len = max_payload;
                if (i>= (int)data_len - (int)max_payload)
                {
                    SEQ |= 0x8000;
                    payload_len = (UInt16)(data_len - i);
                }
                payload_len_all += payload_len;
                Byte[] data_buf2 = data_buffer.Skip(i).ToArray();
                /*Byte[] data_buf2 = new byte[payload_len];//data_buffer.Skip(i).ToArray();
                for(int j=0;j< payload_len; j++)
                {
                    data_buf2[j] = data_buf[i + j];
                }*/
                api_buf = metaverseProtocalGen(version, data_type, cmd_type, ENC, cmd_set, cmd_id, 0, SEQ, data_buf2,payload_len);
                buffer.AddRange(api_buf);
                /*for (int j = 0; j < api_buf.Length; j++)
                {
                    api_buf_all[api_buf_all_idx ++ ] = api_buf[j];
                }*/
                SEQ++;
            }
            //return api_buf_all;
            return buffer.ToArray();
        }
        public Byte[] metaverseProtocalGen(Byte version,Byte data_type,Byte cmd_type,Byte ENC, UInt16 cmd_set, UInt16 cmd_id, Byte extend_len,UInt16 SEQ,  Byte[] data_buf, UInt16 data_len)
        {
            Byte[] api_buf;//= new Byte[];
            UInt16 packet_len_all = (UInt16)(6 + data_len);
            if (version == 0)
            {
                packet_len_all = (UInt16)(4 + data_len);
                if (data_len > meta_api_buf_v0_max_size - 4)
                {
                    packet_len_all = meta_api_buf_v0_max_size;
                    data_len = (UInt16)(packet_len_all - 4);
                }
            }
            if (version == 1)
            {
                packet_len_all = (UInt16)(6 + data_len);
                if (data_len > meta_api_buf_v1_max_size - 6)
                {
                    packet_len_all = meta_api_buf_v1_max_size;
                    data_len = (UInt16)(packet_len_all - 6);
                }
            }
            if (version == 2)
            {
                packet_len_all = (UInt16)(13 + data_len);
                if (data_len > meta_api_buf_v2_max_size - 13)
                {
                    packet_len_all = meta_api_buf_v2_max_size;
                    data_len = (UInt16)(packet_len_all - 13);
                }
            }
            if (version >= 16)
            {
                packet_len_all = (UInt16)(22 + data_len + extend_len);
                if (data_len > meta_api_buf_v16_max_size - 22 - extend_len)
                {
                    packet_len_all = meta_api_buf_v16_max_size;
                    data_len = (UInt16)(packet_len_all - 22 - extend_len);
                }
            }

            if (version < 16)
            {
                UInt16 packet_len = 0;
                
                api_buf = new Byte[packet_len_all];

                if((cmd_type & 0x20) == 0)
                    api_buf[packet_len++] = 0xFA;
                else
                    api_buf[packet_len++] = 0xFB;

                if (version < 2)
                    api_buf[packet_len++] = (Byte)(((version << 6) & 0xFF) | (packet_len_all & 0x3F));
                else
                {
                    api_buf[packet_len++] = (Byte)(((version << 6) & 0xF0) | ((packet_len_all >> 8) & 0x0F));
                    api_buf[packet_len++] = (Byte)(packet_len_all & 0xFF);
                }
                if (version == 1)
                {
                    api_buf[packet_len++] = (Byte)(cmd_set & 0xFF);
                    api_buf[packet_len++] = (Byte)(cmd_id & 0xFF);
                }
                else if(version >= 2)
                {
                    api_buf[packet_len++] = (Byte)(cmd_set & 0xFF);
                    api_buf[packet_len++] = (Byte)((cmd_set >> 8) & 0xFF);
                    api_buf[packet_len++] = (Byte)(cmd_id & 0xFF);
                    api_buf[packet_len++] = (Byte)((cmd_id >> 8) & 0xFF);
                    api_buf[packet_len++] = (Byte)(SEQ & 0xFF);
                    api_buf[packet_len++] = (Byte)((SEQ >> 8) & 0xFF);
                }

                for (int i = 0; i < data_len; i++)
                {
                    api_buf[packet_len++] = data_buf[i];
                }
                if (version <= 1)
                {
                    UInt16 crc16_t = crc16_obj.crc16_init();
                    crc16_t = crc16_obj.crc16_update(crc16_t, api_buf, packet_len);
                    api_buf[packet_len++] = (Byte)(crc16_t & 0xFF);
                    api_buf[packet_len++] = (Byte)((crc16_t >> 8) & 0xFF);
                }
                
                if (version >= 2)
                {
                    UInt32 crc32_t = crc32_obj.crc32_init();
                    crc32_t = crc32_obj.crc32_update(crc32_t, api_buf, packet_len);
                    api_buf[packet_len++] = (Byte)(crc32_t & 0xFF);
                    api_buf[packet_len++] = (Byte)((crc32_t >> 8) & 0xFF);
                    api_buf[packet_len++] = (Byte)((crc32_t >> 16) & 0xFF);
                    api_buf[packet_len++] = (Byte)((crc32_t >> 24) & 0xFF);
                }

            }
            else
            {
                UInt16 packet_len = 0;

                api_buf = new Byte[packet_len_all];
                
                api_buf[packet_len++] = 0x5A;
                api_buf[packet_len++] = 0xEA;
                api_buf[packet_len++] = 0x10;//ver
                api_buf[packet_len++] = data_type;
                api_buf[packet_len++] = (Byte)(packet_len_all & 0xFF);
                api_buf[packet_len++] = (Byte)((packet_len_all >>8) & 0xFF);
                api_buf[packet_len++] = cmd_type;
                api_buf[packet_len++] = ENC;
                api_buf[packet_len++] = (Byte)(cmd_set & 0xFF); 
                api_buf[packet_len++] = (Byte)((cmd_set >> 8) & 0xFF);
                api_buf[packet_len++] = (Byte)(cmd_id & 0xFF);
                api_buf[packet_len++] = (Byte)((cmd_id >> 8) & 0xFF);
                api_buf[packet_len++] = 0;
                api_buf[packet_len++] = 0;
                api_buf[packet_len++] = (Byte)(SEQ & 0xFF);
                api_buf[packet_len++] = (Byte)((SEQ >> 8) & 0xFF);
                UInt16 crc16_t = crc16_obj.crc16_init();
                crc16_t = crc16_obj.crc16_update(crc16_t, api_buf, packet_len);
                api_buf[packet_len++] = (Byte)(crc16_t & 0xFF);
                api_buf[packet_len++] = (Byte)((crc16_t >> 8) & 0xFF);

                for(int i = 0; i < data_len; i++)
                {
                    api_buf[packet_len++] = data_buf[i];
                }

                UInt32 crc32_t = crc32_obj.crc32_init();
                crc32_t = crc32_obj.crc32_update(crc32_t, api_buf, packet_len);
                api_buf[packet_len++] = (Byte)(crc32_t & 0xFF);
                api_buf[packet_len++] = (Byte)((crc32_t >> 8) & 0xFF);
                api_buf[packet_len++] = (Byte)((crc32_t >> 16) & 0xFF);
                api_buf[packet_len++] = (Byte)((crc32_t >> 24) & 0xFF);

            }
            return api_buf;
        }

        private UInt16 get_data_len(Int16 version, UInt16 package_len, UInt16 extend_len)
        {
            UInt16 data_len = 0;
            if (version == 0)
                data_len = (UInt16)(package_len - 4);
            else if (version == 1)
                data_len = (UInt16)(package_len - 6);
            else if (version == 2)
                data_len = (UInt16)(package_len - 13);
            else if (version >= 16)
            {
                data_len = (UInt16)(package_len - 22 - extend_len);
            }
            return data_len;
        }

        private bool meata_version_check(byte version)
        {

            for (int j=0;j< meta_api_support_versions.Length; j++)
            {
                if (version == meta_api_support_versions[j])
                {
                    return true;
                }
            }
            return false;
            /*
            //version check
            if ((meta_api_version < meta_api_min_support_version) ||
                (meta_api_version > meta_api_max_support_version) ||
                (meta_api_version >= 3 && meta_api_version < 16) || meta_api_version > 16)
            {
                meta_api_state = 0;
                break;
            }*/
        }
        #region 协议拆包
        public void metaverseProtocalRecev(Byte[] buf, int size)
        {
            DateTime dt = DateTime.Now;
            //meta_api_crc32_calc = crc32_obj.crc32_init();

            //meta_api_crc32_calc = crc32_obj.crc32_update(meta_api_crc32_calc, buf, 1009);

            //string date1 = dt.ToUniversalTime().ToString("yyyy-MM-dd HH:mm:ss.ffffff");
            try { 
            for (int i = 0; i < size; i++)
            {
                Byte rx_data = buf[i];
                if (meta_api_state >= 1)
                {
                    meta_api_buf[meta_api_buf_idx++] = rx_data;
                    if (meta_api_buf_idx >= meta_api_buf_max_size)
                    {
                        meta_api_state = 0;
                    }
                }
                    switch (meta_api_state)
                    {
                        case 0:
                            meta_api_buf_idx = 0;
                            if (rx_data == 0xFA || rx_data == 0xFB || rx_data == 0x5A)
                            {
                                if (rx_data == 0x5A)
                                {
                                    meta_api_state = 1;
                                }
                                else
                                {
                                    if(meta_api_min_support_version > 2)
                                    {
                                        break;
                                    }
                                    meta_api_state = 2;
                                }
                                meta_api_extend_idx = 0;
                                meta_api_len = 0;
                                meta_api_data_len = 0;
                                
                                meta_api_buf = new Byte[meta_api_buf_max_size];
                                meta_api_buf[meta_api_buf_idx++] = rx_data;
                            }
                            break;
                        case 1:
                            if (rx_data == 0xEA)
                            {
                                meta_api_state = 2;
                            }
                            else
                            {
                                meta_api_state = 0;
                            }
                            break;
                        case 2:

                            if (meta_api_buf[meta_api_buf_idx - 2] == 0xEA)
                            {
                                meta_api_version = (Int16)(rx_data & 0x3F);
                                if (meta_api_version < 16)
                                {
                                    meta_api_state = 0;
                                    break;
                                }
                            }
                            else
                            {
                                meta_api_version = (Byte)(rx_data >> 6);
                            }
                            //version check
                            if (! meata_version_check((Byte)meta_api_version))
                            {
                                meta_api_state = 0;
                                break;
                            }

                            meta_api_rx_idx = 0;

                            if (meta_api_version < 2)
                            {
                                meta_api_len = (Byte)(rx_data & 0x3F);
                                if (meta_api_version == 0)
                                {
                                    meta_api_data_len = get_data_len(meta_api_version, meta_api_len, 0);
                                    if (meta_api_data_len > 0)
                                    {
                                        meta_api_state = 17;
                                    }
                                    else
                                        meta_api_state = 18;
                                }
                                else
                                    meta_api_state = 6;
                            } else if (meta_api_version == 2)
                            {

                                meta_api_len = (UInt16)((rx_data & 0x3F) << 8);
                                meta_api_state = 5;
                            }
                            else
                            {
                                meta_api_state = 3;
                            }
                            break;
                        case 3:
                            {
                                meta_api_data_type = rx_data;
                                meta_api_state = 4;
                            }
                            break;
                        case 4:
                            {
                                meta_api_len = (UInt16)(rx_data);
                                meta_api_state = 5;
                            }
                            break;
                        case 5:
                            {
                                if (meta_api_version == 2)
                                {
                                    meta_api_len |= rx_data;
                                    meta_api_state = 6;
                                }
                                else
                                {
                                    meta_api_len |= (UInt16)(((UInt16)rx_data) << 8);
                                    meta_api_state = 40;
                                }
                            }
                            break;
                        case 40:
                            meta_api_cmd_type = rx_data;
                            meta_api_state = 6;
                            break;
                        case 6:
                            meta_api_cmdset = rx_data;

                            if (meta_api_version < 2)
                                meta_api_state = 8;
                            else
                                meta_api_state = 7;

                            break;
                        case 7:
                            meta_api_cmdset |= (UInt16)(((UInt16)rx_data) << 8);
                            meta_api_state = 8;
                            break;
                        case 8:
                            meta_api_cmdid = rx_data;

                            if (meta_api_version <= 1)
                            {
                                meta_api_data_len = get_data_len(meta_api_version, meta_api_len, 0);
                                if (meta_api_data_len > 0)
                                {
                                    meta_api_state = 17;
                                }
                                else
                                    meta_api_state = 18;
                            }
                            else
                                meta_api_state = 9;
                            break;
                        case 9:
                            meta_api_cmdid |= (UInt16)(((UInt16)rx_data) << 8);

                            if (meta_api_version == 2)
                                meta_api_state = 13;
                            else
                                meta_api_state = 41;
                            break;

                        case 41:
                            //reserv
                            meta_api_state = 10;
                            break;
                        case 10:
                            //reserv
                            meta_api_state = 11;
                            break;
                        case 11:
                            meta_api_extend_len = (Byte)(rx_data & 0x0F);
                            if (meta_api_extend_len > 0)
                                meta_api_state = 12;
                            else
                                meta_api_state = 13;
                            break;
                        case 12:
                            meta_api_extend_idx++;
                            if (meta_api_extend_idx == meta_api_extend_len)
                            {
                                meta_api_state = 13;
                            }
                            break;
                        //seq
                        case 13:
                            meta_api_seq = rx_data;
                            meta_api_state = 14;
                            break;                                                                                                                                                                                                                       //aGV5IHRoaXMgcHJvZ3JhbSBpcyB3cml0dGVuIGJ5IGh0dHBzOi8vZ2l0aHViLmNvbS95b3VrcGFu
                        case 14:
                            meta_api_seq |= (UInt16)(((UInt16)rx_data) << 8);
                            meta_api_state = 15;

                            if (meta_api_version < 16)
                            {
                                meta_api_data_len = get_data_len(meta_api_version, meta_api_len, 0);
                                if (meta_api_data_len > 0)
                                {
                                    meta_api_state = 17;
                                }else
                                    meta_api_state = 18;

                            }
                            else
                                meta_api_state = 15;
                            break;

                        case 15:
                            meta_api_crc16 = (UInt16)rx_data;
                            meta_api_state = 16;
                            break;
                        case 16:
                            meta_api_crc16 |= (UInt16)((UInt16)rx_data << 8);

                            UInt16 crc16_d = crc16_obj.crc16_init();

                            meta_api_crc16_calc = crc16_obj.crc16_update(crc16_d, meta_api_buf, (UInt16)(meta_api_buf_idx - 2));

                            if (meta_api_crc16 != meta_api_crc16_calc)
                            {
                                meta_api_state = 0;
                            }
                            else
                            {
                                meta_api_data_len = get_data_len(meta_api_version, meta_api_len, meta_api_extend_len);
                                if (meta_api_data_len > 0)
                                {
                                    meta_api_state = 17;
                                }
                                else
                                    meta_api_state = 18;
                            }
                            break;
                        case 17:
                            {
                                /*if (meta_api_buf_idx + meta_api_data_len > 65536)
                                {
                                    meta_api_state = 0;
                                    break;
                                }*/
                                //meta_api_state = 43;
                                if (true)
                                {
                                    meta_api_data_buf[meta_api_rx_idx++] = rx_data;
                                    if (meta_api_rx_idx == meta_api_data_len)
                                    {
                                        meta_api_state = 18;
                                        break;
                                    }
                                }else
                                {
                                    //copy meta_api_data_len bytes,for speed up 
                                    //TODO : not pass the test,
                                    
                                    meta_api_data_buf[meta_api_rx_idx++] = rx_data;
                                    if(meta_api_rx_idx!= meta_api_data_len)
                                    {
                                        int k = i;
                                        for (; k < size; k++)
                                        {
                                            rx_data = buf[k];
                                            meta_api_data_buf[meta_api_rx_idx++] = rx_data;
                                            meta_api_buf[meta_api_buf_idx++] = rx_data;
                                            if (meta_api_rx_idx == meta_api_data_len)
                                            {
                                                break;
                                            }

                                        }
                                        i =  k;
                                    }
                                    else
                                    {
                                        break;
                                    }

                                    meta_api_state = 18;
                                }

                                break;
                            }
                    case 18:
                        if (meta_api_version <=1)
                            meta_api_crc16 = (UInt16)rx_data;
                        else
                            meta_api_crc32 = (UInt32)rx_data;

                        meta_api_state = 19;
                        break;
                    case 19:

                        if (meta_api_version <= 1)
                        {
                            meta_api_crc16 |= (UInt16)((UInt16)rx_data << 8);
                            meta_api_crc16_calc = crc16_obj.crc16_init();

                            meta_api_crc16_calc = crc16_obj.crc16_update(meta_api_crc16_calc, meta_api_buf, (UInt16)(meta_api_buf_idx - 2));

                            if (meta_api_crc16 == meta_api_crc16_calc)
                            {
                                Process_metaverse_packet((Byte)meta_api_version, meta_api_cmdset, meta_api_cmdid, meta_api_seq, meta_api_buf, meta_api_len, meta_api_data_buf, meta_api_rx_idx);
                            }

                            meta_api_state = 0;
                        }
                        else
                        {
                            meta_api_crc32 |= (UInt32)((UInt32)rx_data << 8);
                            meta_api_state = 20;
                        }

                        break;

                    case 20:
                        meta_api_crc32 |= (UInt32)((UInt32)rx_data << 16);
                        meta_api_state = 21;
                        break;
                    case 21:
                        meta_api_crc32 |= (UInt32)((UInt32)rx_data << 24);
                        meta_api_crc32_calc = crc32_obj.crc32_init();
                        meta_api_crc32_calc = crc32_obj.crc32_update(meta_api_crc32_calc, meta_api_buf, (UInt16)(meta_api_buf_idx - 4));

                        if(meta_api_crc32_calc == meta_api_crc32)
                        {
                            Process_metaverse_packet((Byte)meta_api_version,meta_api_cmdset, meta_api_cmdid, meta_api_seq, meta_api_buf, meta_api_len, meta_api_data_buf, meta_api_rx_idx);
                        }
                        
                        meta_api_state = 0;
                        break;

                }
            }

            }
            catch (System.Exception ex)
            {
                return; 
            }
        }
        #endregion

        private void save_packet(Byte [] package_buf_all)
        {
            long time_passed_ms = DateTime.Now.ToUniversalTime().Ticks / 10000 - mainform.start_time_ms;                                                                                                                                                                      // //im the author :dGhpcyBwcm9ncmFtIGlzIHdyaXR0ZW4gYnkgeW91a3BhbkBnbWFpbC5jb20=

            if (mainform != null)
            {
                mainform.rcv_pkg_counter_all++;
                if (mainform.saveDataFS != null && mainform.radioButtonFileSave.Checked)
                {
                    try
                    {
                        var time_info = "";
                        //version,timstamp,counter,buf_data
                        //var data_hex_str = BitConverter.ToString(package_buf_all, 0, meta_api_len).Replace("-", string.Empty);
                        var data_hex_str = Convert.ToBase64String(package_buf_all);
                        if (last_time_ticks != time_passed_ms)
                        {
                            last_time_ticks = time_passed_ms;
                            time_info = time_passed_ms.ToString();
                        }

                        if (last_save_data != data_hex_str)
                        {
                            last_save_data = data_hex_str;
                            var write_msg = mainform.rcv_pkg_counter_all + "," + time_passed_ms + "," + data_hex_str + "\n";
                            byte[] info = new UTF8Encoding(true).GetBytes(write_msg);
                            mainform.saveDataFS.Write(info, 0, info.Length);
                        }

                    }
                    catch (System.Exception ex)
                    {

                    }
                }
            }
        }

        public void unit_test()
        {
            //(Byte version, Byte data_type, Byte cmd_type, Byte ENC, UInt16 cmd_set, UInt16 cmd_id, Byte extend_len, UInt16 SEQ, Byte[] data_buf, UInt16 data_len)
            UInt32 data_len = 120000;
            var data_buf = new Byte[data_len];
            Random rd = new Random();
            for (int i = 0; i < data_len; i++)
            {
                data_buf[i] = (Byte)i;
            }
            Byte version;
            Byte data_type = 1;
            Byte cmd_type = 0x40;
            Byte ENC = 0;
            UInt16 cmd_set = 0;
            UInt16 cmd_id = 0;
            Byte extend_len = 0;
            UInt16 SEQ = 0x8000;

            Byte[] version_test = new byte[] { 0, 1, 2, 16 };
            long start_time_pc_ms = DateTime.Now.ToUniversalTime().Ticks / 10000;
            
            for (int i = 0; i < version_test.Length; i++)
            {
                version = version_test[i];
                var data_len2 = data_len;
                if (data_len2 > 65535 - 22)
                {
                    data_len2 = 65535 - 22;
                }
                var buf = metaverseProtocalGen(version, data_type, cmd_type, ENC, cmd_set, cmd_id, extend_len, SEQ, data_buf, (UInt16)data_len2);
                //for test err
                if (false)
                {
                    int rand = rd.Next(100);
                    if (rand > 80)
                    {
                        mainform.msg_box_print(String.Format(@"insert err byte in version {0:D}  ", version));
                        buf[rd.Next(buf.Length - 1)] = (Byte)rd.Next(255);
                    }
                }

                metaverseProtocalRecev(buf, buf.Length);
            }
            //not recommend version <=2 to gen massive buf
            version_test = new byte[] {   2,  16 };
            for (int i = 0; i < version_test.Length; i++)
            {
                version = version_test[i];

                var buf1 = metaverseProtocalGenMassive(version, data_type, cmd_type, ENC, cmd_set, cmd_id, data_buf, data_len);
                if (false)
                {
                    int err_bytes = 3;
                    for (int j = 0; j < err_bytes; j++)
                    {
                        int err_idx = rd.Next(buf1.Length - 1);
                        //mainform.msg_box_print(String.Format(@"insert err byte in version {0:D} ,err_idx: {1:D} ", version, err_idx));
                        buf1[err_idx] = (Byte)rd.Next(255);
                    }

                    mainform.msg_box_print(String.Format(@"insert err {0:D} byte in version {1:D} ",err_bytes,version));
                }
                metaverseProtocalRecev(buf1, buf1.Length);
            }

            var buf2 = metaverseProtocalGenMassive(16, data_type, cmd_type, ENC, cmd_set, cmd_id, data_buf, data_len);

            metaverseProtocalRecev(buf2, buf2.Length);

            //the print will need time ,so will clear quick ,to test the real time use
            mainform.msg_box_print(String.Format(@"test finished,send&recv time used  {0:F3} ms ", DateTime.Now.ToUniversalTime().Ticks / 10000 - start_time_pc_ms));
            if(mainform.textBoxReceive.Text.Length > 20000)
            {
                mainform.textBoxReceive.Text = "";
            }
        }
    }
}
