package cc

var hdr_u_h = `
typedef signed char schar;
typedef unsigned char uchar;
typedef unsigned short ushort;
typedef unsigned int uint;
typedef unsigned long ulong;
typedef long long vlong;
typedef unsigned long long uvlong;

typedef schar int8;
typedef uchar uint8;
typedef short int16;
typedef ushort uint16;
typedef long int32;
typedef ulong uint32;
typedef vlong int64;
typedef uvlong uint64;

typedef schar s8int;
typedef uchar u8int;
typedef short s16int;
typedef ushort u16int;
typedef long s32int;
typedef ulong u32int;
typedef vlong s64int;
typedef uvlong u64int;

typedef uint32 Rune;

typedef struct va_list *va_list;
`

var hdr_libc_h = `

int memcmp(void*, void*, long);
void *memset(void*, int, long);
int strcmp(char*, char*);
int strncmp(char*, char*, int);
char *strcpy(char*, char*);

int errstr(char*, uint);
void werrstr(char*, ...);

int fprint(int, char*, ...);
int snprint(char*, int, char*, ...);
char *seprint(char*, char*, char*, ...);
char *vseprint(char*, char*, char*, va_list);

void va_start(va_list, void*);
void va_end(va_list);
`
