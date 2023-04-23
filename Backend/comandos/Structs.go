package comandos

type Partition struct {
	Part_status uint8     // Cambiar char a uint8
	Part_type   uint8     // Cambiar char a uint8
	Part_fit    uint8     // Cambiar char a uint8
	Part_start  int64     // Cambiar int a int64
	Part_s      int64     // Cambiar int a int64
	Part_name   [16]uint8 // Cambiar char a uint8
}

type MBR struct {
	Mbr_tamano         int64 // Cambiar int a int64
	Mbr_fecha_creacion int64 // Cambiar time_t a int64
	Mbr_dsk_signature  int64 // Cambiar int a int64
	Disk_fit           uint8 // Cambiar char a uint8
	Mbr_partition      [4]Partition
}

type EBR struct {
	Part_status uint8     // Cambiar char a uint8
	Part_fit    uint8     // Cambiar char a uint8
	Part_start  int64     // Cambiar int a int64
	Part_s      int64     // Cambiar int a int64
	Part_next   int64     // Cambiar int a int64
	Part_name   [16]uint8 // Cambiar char a uint8
}

type SuperBloque struct {
	S_filesystem_type   int64 // Cambiar int a int64
	S_inodes_count      int64 // Cambiar int a int64
	S_blocks_count      int64 // Cambiar int a int64
	S_free_blocks_count int64 // Cambiar int a int64
	S_free_inodes_count int64 // Cambiar int a int64
	S_mtime             int64 // Cambiar time_t a int64
	S_umtime            int64 // Cambiar time_t a int64
	S_mnt_count         int64 // Cambiar int a int64
	S_magic             int64 // Cambiar int a int64
	S_inode_s           int64 // Cambiar int a int64
	S_block_s           int64 // Cambiar int a int64
	S_firts_ino         int64 // Cambiar int a int64
	S_first_blo         int64 // Cambiar int a int64
	S_bm_inode_start    int64 // Cambiar int a int64
	S_bm_block_start    int64 // Cambiar int a int64
	S_inode_start       int64 // Cambiar int a int64
	S_block_start       int64 // Cambiar int a int64
}

type TablaInodo struct {
	I_uid   int64     // Cambiar int a int64
	I_gid   int64     // Cambiar int a int64
	I_s     int64     // Cambiar int a int64
	I_atime int64     // Cambiar time_t a int64
	I_ctime int64     // Cambiar time_t a int64
	I_mtime int64     // Cambiar time_t a int64
	I_block [15]int64 // Cambiar int a int64
	I_type  uint8     // Cambiar char a uint8
	I_perm  int64     // Cambiar int a int64
}

type Content struct {
	B_name  [12]uint8 // Cambiar char a uint8
	B_inodo int64     // Cambiar int a int64
}

type BloqueCarpeta struct {
	B_content [4]Content
}

type BloqueArchivo struct {
	B_content [64]uint8 // Cambiar char a uint8
}

type BloqueApuntador struct {
	B_pointers [16]int64 // Cambiar int a int64
}

type Journal struct {
	Journal_Tipo_Operacion [10]uint8  // Cambiar char a uint8
	Journal_Tipo           uint8      // Cambiar char a uint8
	Journal_Path           [100]uint8 // Cambiar char a uint8
	Journal_Contenido      [100]uint8 // Cambiar char a uint8
	Journal_Fecha          int64      // Cambiar time_t a int64
	Journal_Size           int64      // Cambiar int a int64
	Journal_Sig            int64      // Cambiar int a int64
	Journal_Start          int64      // Cambiar int a int64
}
