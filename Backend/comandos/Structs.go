package comandos

import "time"

type Partition struct {
	part_status byte
	part_type   byte
	part_fit    byte
	part_start  int32
	part_s      int32
	part_name   [16]byte
}

type MBR struct {
	mbr_tamano         int32
	mbr_fecha_creacion time.Time // recordar que al convertir el tiempo ejemplo: fechaFormateada := mbr.mbr_fecha_creacion.Format("2006-01-02 15:04:05")
	mbr_dsk_signature  int32
	disk_fit           byte
	mbr_partition      [4]Partition
}

type EBR struct {
	part_status byte
	part_fit    byte
	part_start  int32
	part_s      int32
	part_next   int32
	part_name   [16]byte
}

type SuperBloque struct {
	s_filesystem_type   int32
	s_inodes_count      int32
	s_blocks_count      int32
	s_free_blocks_count int32
	s_free_inodes_count int32
	s_mtime             time.Time
	s_umtime            time.Time
	s_mnt_count         int32
	s_magic             int32
	s_inode_s           int32
	s_block_s           int32
	s_firts_ino         int32
	s_first_blo         int32
	s_bm_inode_start    int32
	s_bm_block_start    int32
	s_inode_start       int32
	s_block_start       int32
}

type TablaInodo struct {
	i_uid   int32
	i_gid   int32
	i_size  int32
	i_atime time.Time
	i_ctime time.Time
	i_mtime time.Time
	i_block [15]int32
	i_type  byte
	i_perm  int32
}

type content struct {
	b_name  [12]byte
	b_inodo int32
}

type BloqueCarpeta struct {
	b_content [4]content
}

type BloqueArchivo struct {
	b_content [64]byte
}

type BloqueApuntador struct {
	b_pointers [16]int32
}

type Journal struct {
	journal_Tipo_Operacion [10]byte
	journal_Tipo           byte
	journal_Path           [100]byte
	journal_Contenido      [100]byte
	journal_Fecha          time.Time
	journal_Size           int32
	journal_Sig            int32
	journal_Start          int32
}
