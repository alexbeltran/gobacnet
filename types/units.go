package types

type Unit uint32

const (
	/* Acceleration */
	UNITS_METERS_PER_SECOND_PER_SECOND Unit = 166
	/* Area */
	UNITS_SQUARE_METERS      Unit = 0
	UNITS_SQUARE_CENTIMETERS Unit = 116
	UNITS_SQUARE_FEET        Unit = 1
	UNITS_SQUARE_INCHES      Unit = 115
	/* Currency */
	UNITS_CURRENCY1  Unit = 105
	UNITS_CURRENCY2  Unit = 106
	UNITS_CURRENCY3  Unit = 107
	UNITS_CURRENCY4  Unit = 108
	UNITS_CURRENCY5  Unit = 109
	UNITS_CURRENCY6  Unit = 110
	UNITS_CURRENCY7  Unit = 111
	UNITS_CURRENCY8  Unit = 112
	UNITS_CURRENCY9  Unit = 113
	UNITS_CURRENCY10 Unit = 114
	/* Electrical */
	UNITS_MILLIAMPERES              Unit = 2
	UNITS_AMPERES                   Unit = 3
	UNITS_AMPERES_PER_METER         Unit = 167
	UNITS_AMPERES_PER_SQUARE_METER  Unit = 168
	UNITS_AMPERE_SQUARE_METERS      Unit = 169
	UNITS_DECIBELS                  Unit = 199
	UNITS_DECIBELS_MILLIVOLT        Unit = 200
	UNITS_DECIBELS_VOLT             Unit = 201
	UNITS_FARADS                    Unit = 170
	UNITS_HENRYS                    Unit = 171
	UNITS_OHMS                      Unit = 4
	UNITS_OHM_METERS                Unit = 172
	UNITS_MILLIOHMS                 Unit = 145
	UNITS_KILOHMS                   Unit = 122
	UNITS_MEGOHMS                   Unit = 123
	UNITS_MICROSIEMENS              Unit = 190
	UNITS_MILLISIEMENS              Unit = 202
	UNITS_SIEMENS                   Unit = 173 /* 1 mho equals 1 siemens */
	UNITS_SIEMENS_PER_METER         Unit = 174
	UNITS_TESLAS                    Unit = 175
	UNITS_VOLTS                     Unit = 5
	UNITS_MILLIVOLTS                Unit = 124
	UNITS_KILOVOLTS                 Unit = 6
	UNITS_MEGAVOLTS                 Unit = 7
	UNITS_VOLT_AMPERES              Unit = 8
	UNITS_KILOVOLT_AMPERES          Unit = 9
	UNITS_MEGAVOLT_AMPERES          Unit = 10
	UNITS_VOLT_AMPERES_REACTIVE     Unit = 11
	UNITS_KILOVOLT_AMPERES_REACTIVE Unit = 12
	UNITS_MEGAVOLT_AMPERES_REACTIVE Unit = 13
	UNITS_VOLTS_PER_DEGREE_KELVIN   Unit = 176
	UNITS_VOLTS_PER_METER           Unit = 177
	UNITS_DEGREES_PHASE             Unit = 14
	UNITS_POWER_FACTOR              Unit = 15
	UNITS_WEBERS                    Unit = 178
	/* Energy */
	UNITS_JOULES                  Unit = 16
	UNITS_KILOJOULES              Unit = 17
	UNITS_KILOJOULES_PER_KILOGRAM Unit = 125
	UNITS_MEGAJOULES              Unit = 126
	UNITS_WATT_HOURS              Unit = 18
	UNITS_KILOWATT_HOURS          Unit = 19
	UNITS_MEGAWATT_HOURS          Unit = 146
	UNITS_WATT_HOURS_REACTIVE     Unit = 203
	UNITS_KILOWATT_HOURS_REACTIVE Unit = 204
	UNITS_MEGAWATT_HOURS_REACTIVE Unit = 205
	UNITS_BTUS                    Unit = 20
	UNITS_KILO_BTUS               Unit = 147
	UNITS_MEGA_BTUS               Unit = 148
	UNITS_THERMS                  Unit = 21
	UNITS_TON_HOURS               Unit = 22
	/* Enthalpy */
	UNITS_JOULES_PER_KILOGRAM_DRY_AIR     Unit = 23
	UNITS_KILOJOULES_PER_KILOGRAM_DRY_AIR Unit = 149
	UNITS_MEGAJOULES_PER_KILOGRAM_DRY_AIR Unit = 150
	UNITS_BTUS_PER_POUND_DRY_AIR          Unit = 24
	UNITS_BTUS_PER_POUND                  Unit = 117
	/* Entropy */
	UNITS_JOULES_PER_DEGREE_KELVIN          Unit = 127
	UNITS_KILOJOULES_PER_DEGREE_KELVIN      Unit = 151
	UNITS_MEGAJOULES_PER_DEGREE_KELVIN      Unit = 152
	UNITS_JOULES_PER_KILOGRAM_DEGREE_KELVIN Unit = 128
	/* Force */
	UNITS_NEWTON Unit = 153
	/* Frequency */
	UNITS_CYCLES_PER_HOUR   Unit = 25
	UNITS_CYCLES_PER_MINUTE Unit = 26
	UNITS_HERTZ             Unit = 27
	UNITS_KILOHERTZ         Unit = 129
	UNITS_MEGAHERTZ         Unit = 130
	UNITS_PER_HOUR          Unit = 131
	/* Humidity */
	UNITS_GRAMS_OF_WATER_PER_KILOGRAM_DRY_AIR Unit = 28
	UNITS_PERCENT_RELATIVE_HUMIDITY           Unit = 29
	/* Length */
	UNITS_MICROMETERS Unit = 194
	UNITS_MILLIMETERS Unit = 30
	UNITS_CENTIMETERS Unit = 118
	UNITS_KILOMETERS  Unit = 193
	UNITS_METERS      Unit = 31
	UNITS_INCHES      Unit = 32
	UNITS_FEET        Unit = 33
	/* Light */
	UNITS_CANDELAS                  Unit = 179
	UNITS_CANDELAS_PER_SQUARE_METER Unit = 180
	UNITS_WATTS_PER_SQUARE_FOOT     Unit = 34
	UNITS_WATTS_PER_SQUARE_METER    Unit = 35
	UNITS_LUMENS                    Unit = 36
	UNITS_LUXES                     Unit = 37
	UNITS_FOOT_CANDLES              Unit = 38
	/* Mass */
	UNITS_MILLIGRAMS  Unit = 196
	UNITS_GRAMS       Unit = 195
	UNITS_KILOGRAMS   Unit = 39
	UNITS_POUNDS_MASS Unit = 40
	UNITS_TONS        Unit = 41
	/* Mass Flow */
	UNITS_GRAMS_PER_SECOND       Unit = 154
	UNITS_GRAMS_PER_MINUTE       Unit = 155
	UNITS_KILOGRAMS_PER_SECOND   Unit = 42
	UNITS_KILOGRAMS_PER_MINUTE   Unit = 43
	UNITS_KILOGRAMS_PER_HOUR     Unit = 44
	UNITS_POUNDS_MASS_PER_SECOND Unit = 119
	UNITS_POUNDS_MASS_PER_MINUTE Unit = 45
	UNITS_POUNDS_MASS_PER_HOUR   Unit = 46
	UNITS_TONS_PER_HOUR          Unit = 156
	/* Power */
	UNITS_MILLIWATTS         Unit = 132
	UNITS_WATTS              Unit = 47
	UNITS_KILOWATTS          Unit = 48
	UNITS_MEGAWATTS          Unit = 49
	UNITS_BTUS_PER_HOUR      Unit = 50
	UNITS_KILO_BTUS_PER_HOUR Unit = 157
	UNITS_HORSEPOWER         Unit = 51
	UNITS_TONS_REFRIGERATION Unit = 52
	/* Pressure */
	UNITS_PASCALS                      Unit = 53
	UNITS_HECTOPASCALS                 Unit = 133
	UNITS_KILOPASCALS                  Unit = 54
	UNITS_MILLIBARS                    Unit = 134
	UNITS_BARS                         Unit = 55
	UNITS_POUNDS_FORCE_PER_SQUARE_INCH Unit = 56
	UNITS_MILLIMETERS_OF_WATER         Unit = 206
	UNITS_CENTIMETERS_OF_WATER         Unit = 57
	UNITS_INCHES_OF_WATER              Unit = 58
	UNITS_MILLIMETERS_OF_MERCURY       Unit = 59
	UNITS_CENTIMETERS_OF_MERCURY       Unit = 60
	UNITS_INCHES_OF_MERCURY            Unit = 61
	/* Temperature */
	UNITS_DEGREES_CELSIUS           Unit = 62
	UNITS_DEGREES_KELVIN            Unit = 63
	UNITS_DEGREES_KELVIN_PER_HOUR   Unit = 181
	UNITS_DEGREES_KELVIN_PER_MINUTE Unit = 182
	UNITS_DEGREES_FAHRENHEIT        Unit = 64
	UNITS_DEGREE_DAYS_CELSIUS       Unit = 65
	UNITS_DEGREE_DAYS_FAHRENHEIT    Unit = 66
	UNITS_DELTA_DEGREES_FAHRENHEIT  Unit = 120
	UNITS_DELTA_DEGREES_KELVIN      Unit = 121
	/* Time */
	UNITS_YEARS              Unit = 67
	UNITS_MONTHS             Unit = 68
	UNITS_WEEKS              Unit = 69
	UNITS_DAYS               Unit = 70
	UNITS_HOURS              Unit = 71
	UNITS_MINUTES            Unit = 72
	UNITS_SECONDS            Unit = 73
	UNITS_HUNDREDTHS_SECONDS Unit = 158
	UNITS_MILLISECONDS       Unit = 159
	/* Torque */
	UNITS_NEWTON_METERS Unit = 160
	/* Velocity */
	UNITS_MILLIMETERS_PER_SECOND Unit = 161
	UNITS_MILLIMETERS_PER_MINUTE Unit = 162
	UNITS_METERS_PER_SECOND      Unit = 74
	UNITS_METERS_PER_MINUTE      Unit = 163
	UNITS_METERS_PER_HOUR        Unit = 164
	UNITS_KILOMETERS_PER_HOUR    Unit = 75
	UNITS_FEET_PER_SECOND        Unit = 76
	UNITS_FEET_PER_MINUTE        Unit = 77
	UNITS_MILES_PER_HOUR         Unit = 78
	/* Volume */
	UNITS_CUBIC_FEET       Unit = 79
	UNITS_CUBIC_METERS     Unit = 80
	UNITS_IMPERIAL_GALLONS Unit = 81
	UNITS_MILLILITERS      Unit = 197
	UNITS_LITERS           Unit = 82
	UNITS_US_GALLONS       Unit = 83
	/* Volumetric Flow */
	UNITS_CUBIC_FEET_PER_SECOND       Unit = 142
	UNITS_CUBIC_FEET_PER_MINUTE       Unit = 84
	UNITS_CUBIC_FEET_PER_HOUR         Unit = 191
	UNITS_CUBIC_METERS_PER_SECOND     Unit = 85
	UNITS_CUBIC_METERS_PER_MINUTE     Unit = 165
	UNITS_CUBIC_METERS_PER_HOUR       Unit = 135
	UNITS_IMPERIAL_GALLONS_PER_MINUTE Unit = 86
	UNITS_MILLILITERS_PER_SECOND      Unit = 198
	UNITS_LITERS_PER_SECOND           Unit = 87
	UNITS_LITERS_PER_MINUTE           Unit = 88
	UNITS_LITERS_PER_HOUR             Unit = 136
	UNITS_US_GALLONS_PER_MINUTE       Unit = 89
	UNITS_US_GALLONS_PER_HOUR         Unit = 192
	/* Other */
	UNITS_DEGREES_ANGULAR                        Unit = 90
	UNITS_DEGREES_CELSIUS_PER_HOUR               Unit = 91
	UNITS_DEGREES_CELSIUS_PER_MINUTE             Unit = 92
	UNITS_DEGREES_FAHRENHEIT_PER_HOUR            Unit = 93
	UNITS_DEGREES_FAHRENHEIT_PER_MINUTE          Unit = 94
	UNITS_JOULE_SECONDS                          Unit = 183
	UNITS_KILOGRAMS_PER_CUBIC_METER              Unit = 186
	UNITS_KW_HOURS_PER_SQUARE_METER              Unit = 137
	UNITS_KW_HOURS_PER_SQUARE_FOOT               Unit = 138
	UNITS_MEGAJOULES_PER_SQUARE_METER            Unit = 139
	UNITS_MEGAJOULES_PER_SQUARE_FOOT             Unit = 140
	UNITS_NO_UNITS                               Unit = 95
	UNITS_NEWTON_SECONDS                         Unit = 187
	UNITS_NEWTONS_PER_METER                      Unit = 188
	UNITS_PARTS_PER_MILLION                      Unit = 96
	UNITS_PARTS_PER_BILLION                      Unit = 97
	UNITS_PERCENT                                Unit = 98
	UNITS_PERCENT_OBSCURATION_PER_FOOT           Unit = 143
	UNITS_PERCENT_OBSCURATION_PER_METER          Unit = 144
	UNITS_PERCENT_PER_SECOND                     Unit = 99
	UNITS_PER_MINUTE                             Unit = 100
	UNITS_PER_SECOND                             Unit = 101
	UNITS_PSI_PER_DEGREE_FAHRENHEIT              Unit = 102
	UNITS_RADIANS                                Unit = 103
	UNITS_RADIANS_PER_SECOND                     Unit = 184
	UNITS_REVOLUTIONS_PER_MINUTE                 Unit = 104
	UNITS_SQUARE_METERS_PER_NEWTON               Unit = 185
	UNITS_WATTS_PER_METER_PER_DEGREE_KELVIN      Unit = 189
	UNITS_WATTS_PER_SQUARE_METER_DEGREE_KELVIN   Unit = 141
	UNITS_PER_MILLE                              Unit = 207
	UNITS_GRAMS_PER_GRAM                         Unit = 208
	UNITS_KILOGRAMS_PER_KILOGRAM                 Unit = 209
	UNITS_GRAMS_PER_KILOGRAM                     Unit = 210
	UNITS_MILLIGRAMS_PER_GRAM                    Unit = 211
	UNITS_MILLIGRAMS_PER_KILOGRAM                Unit = 212
	UNITS_GRAMS_PER_MILLILITER                   Unit = 213
	UNITS_GRAMS_PER_LITER                        Unit = 214
	UNITS_MILLIGRAMS_PER_LITER                   Unit = 215
	UNITS_MICROGRAMS_PER_LITER                   Unit = 216
	UNITS_GRAMS_PER_CUBIC_METER                  Unit = 217
	UNITS_MILLIGRAMS_PER_CUBIC_METER             Unit = 218
	UNITS_MICROGRAMS_PER_CUBIC_METER             Unit = 219
	UNITS_NANOGRAMS_PER_CUBIC_METER              Unit = 220
	UNITS_GRAMS_PER_CUBIC_CENTIMETER             Unit = 221
	UNITS_BECQUERELS                             Unit = 222
	UNITS_KILOBECQUERELS                         Unit = 223
	UNITS_MEGABECQUERELS                         Unit = 224
	UNITS_GRAY                                   Unit = 225
	UNITS_MILLIGRAY                              Unit = 226
	UNITS_MICROGRAY                              Unit = 227
	UNITS_SIEVERTS                               Unit = 228
	UNITS_MILLISIEVERTS                          Unit = 229
	UNITS_MICROSIEVERTS                          Unit = 230
	UNITS_MICROSIEVERTS_PER_HOUR                 Unit = 231
	UNITS_DECIBELS_A                             Unit = 232
	UNITS_NEPHELOMETRIC_TURBIDITY_UNIT           Unit = 233
	UNITS_PH                                     Unit = 234
	UNITS_GRAMS_PER_SQUARE_METER                 Unit = 235
	UNITS_MINUTES_PER_DEGREE_KELVIN              Unit = 236
	UNITS_OHM_METER_SQUARED_PER_METER            Unit = 237
	UNITS_AMPERE_SECONDS                         Unit = 238
	UNITS_VOLT_AMPERE_HOURS                      Unit = 239
	UNITS_KILOVOLT_AMPERE_HOURS                  Unit = 240
	UNITS_MEGAVOLT_AMPERE_HOURS                  Unit = 241
	UNITS_VOLT_AMPERE_HOURS_REACTIVE             Unit = 242
	UNITS_KILOVOLT_AMPERE_HOURS_REACTIVE         Unit = 243
	UNITS_MEGAVOLT_AMPERE_HOURS_REACTIVE         Unit = 244
	UNITS_VOLT_SQUARE_HOURS                      Unit = 245
	UNITS_AMPERE_SQUARE_HOURS                    Unit = 246
	UNITS_JOULE_PER_HOURS                        Unit = 247
	UNITS_CUBIC_FEET_PER_DAY                     Unit = 248
	UNITS_CUBIC_METERS_PER_DAY                   Unit = 249
	UNITS_WATT_HOURS_PER_CUBIC_METER             Unit = 250
	UNITS_JOULES_PER_CUBIC_METER                 Unit = 251
	UNITS_MOLE_PERCENT                           Unit = 252
	UNITS_PASCAL_SECONDS                         Unit = 253
	UNITS_MILLION_STANDARD_CUBIC_FEET_PER_MINUTE Unit = 254
	/* 255 - NOT USED */
	UNITS_STANDARD_CUBIC_FEET_PER_DAY          Unit = 47808
	UNITS_MILLION_STANDARD_CUBIC_FEET_PER_DAY  Unit = 47809
	UNITS_THOUSAND_CUBIC_FEET_PER_DAY          Unit = 47810
	UNITS_THOUSAND_STANDARD_CUBIC_FEET_PER_DAY Unit = 47811
	UNITS_POUNDS_MASS_PER_DAY                  Unit = 47812
	/* 47813 - NOT USED */
	UNITS_MILLIREMS          Unit = 47814
	UNITS_MILLIREMS_PER_HOUR Unit = 47815
	/* Enumerated values 0-255 and 47808-49999 are reserved for
	   definition by ASHRAE. */
	/* Enumerated values 256-47807 and 50000-65535 may be used by others
	   subject to the procedures and constraints described in Clause 23. */
	/* do the proprietary range inside of enum so that
	   compilers will allocate adequate sized datatype for enum
	   which is used to store decoding */
	UNITS_PROPRIETARY_RANGE_MIN Unit = 256
	UNITS_PROPRIETARY_RANGE_MAX Unit = 65535
)
