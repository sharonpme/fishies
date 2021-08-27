EESchema Schematic File Version 4
EELAYER 30 0
EELAYER END
$Descr A4 11693 8268
encoding utf-8
Sheet 1 1
Title ""
Date ""
Rev ""
Comp ""
Comment1 ""
Comment2 ""
Comment3 ""
Comment4 ""
$EndDescr
$Comp
L Connector:USB_B_Micro J1
U 1 1 6125A1DB
P 3700 4000
F 0 "J1" H 3757 4467 50  0000 C CNN
F 1 "USB" H 3757 4376 50  0000 C CNN
F 2 "Connector_USB:USB_Micro-B_Amphenol_10103594-0001LF_Horizontal" H 3850 3950 50  0001 C CNN
F 3 "~" H 3850 3950 50  0001 C CNN
	1    3700 4000
	1    0    0    -1  
$EndComp
$Comp
L MCU_Module:WeMos_D1_mini U1
U 1 1 6125745A
P 5550 4200
F 0 "U1" H 5550 3311 50  0000 C CNN
F 1 "WeMos_D1_mini" H 5550 3220 50  0000 C CNN
F 2 "Module:WEMOS_D1_mini_light" H 5550 3050 50  0001 C CNN
F 3 "https://wiki.wemos.cc/products:d1:d1_mini#documentation" H 3700 3050 50  0001 C CNN
	1    5550 4200
	1    0    0    -1  
$EndComp
$Comp
L Device:R_Small_US R1
U 1 1 61263487
P 6250 4300
F 0 "R1" V 6045 4300 50  0000 C CNN
F 1 "1k" V 6136 4300 50  0000 C CNN
F 2 "Resistor_SMD:R_0805_2012Metric_Pad1.20x1.40mm_HandSolder" H 6250 4300 50  0001 C CNN
F 3 "~" H 6250 4300 50  0001 C CNN
	1    6250 4300
	0    1    1    0   
$EndComp
$Comp
L Device:LED D1
U 1 1 612648D4
P 6500 4600
F 0 "D1" V 6539 4482 50  0000 R CNN
F 1 "LED" V 6448 4482 50  0000 R CNN
F 2 "Connector_JST:JST_XH_B2B-XH-AM_1x02_P2.50mm_Vertical" H 6500 4600 50  0001 C CNN
F 3 "~" H 6500 4600 50  0001 C CNN
	1    6500 4600
	0    -1   -1   0   
$EndComp
$Comp
L Device:CP C1
U 1 1 612657E0
P 4500 4150
F 0 "C1" H 4618 4196 50  0000 L CNN
F 1 "470uF" H 4618 4105 50  0000 L CNN
F 2 "Capacitor_THT:CP_Radial_D8.0mm_P2.50mm" H 4538 4000 50  0001 C CNN
F 3 "~" H 4500 4150 50  0001 C CNN
	1    4500 4150
	1    0    0    -1  
$EndComp
Wire Wire Line
	3600 4400 3700 4400
Wire Wire Line
	3700 4400 3700 5000
Wire Wire Line
	3700 5000 4500 5000
Connection ~ 3700 4400
Wire Wire Line
	5550 5000 6500 5000
Wire Wire Line
	6500 5000 6500 4750
Connection ~ 5550 5000
Wire Wire Line
	4500 4300 4500 5000
Connection ~ 4500 5000
Wire Wire Line
	4500 5000 5550 5000
Wire Wire Line
	4000 3800 4500 3800
Wire Wire Line
	4500 3800 4500 4000
Wire Wire Line
	4500 3800 4500 3200
Wire Wire Line
	4500 3200 5450 3200
Wire Wire Line
	5450 3200 5450 3400
Connection ~ 4500 3800
Wire Wire Line
	5950 4300 6150 4300
Wire Wire Line
	6350 4300 6500 4300
Wire Wire Line
	6500 4300 6500 4450
$Comp
L Connector:Conn_01x03_Male J2
U 1 1 612672FA
P 7150 4300
F 0 "J2" H 7122 4232 50  0000 R CNN
F 1 "Servo" H 7122 4323 50  0000 R CNN
F 2 "Connector_PinHeader_2.54mm:PinHeader_1x03_P2.54mm_Vertical" H 7150 4300 50  0001 C CNN
F 3 "~" H 7150 4300 50  0001 C CNN
	1    7150 4300
	-1   0    0    1   
$EndComp
Wire Wire Line
	6950 4400 6950 5000
Wire Wire Line
	6950 5000 6500 5000
Connection ~ 6500 5000
Wire Wire Line
	5950 4200 6950 4200
Wire Wire Line
	6950 4300 6700 4300
Wire Wire Line
	6700 4300 6700 3200
Wire Wire Line
	6700 3200 5450 3200
Connection ~ 5450 3200
$EndSCHEMATC
