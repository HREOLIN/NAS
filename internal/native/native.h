#ifndef NAS_NATIVE_H
#define NAS_NATIVE_H

int nas_recovery_create_snapshot(const char* volume, const char* name, char* out, int out_len);
int nas_recovery_restore_file(const char* snapshot_id, const char* source_path, const char* target_path, char* out, int out_len);
int nas_recovery_rollback_volume(const char* snapshot_id, const char* volume, char* out, int out_len);
int nas_network_apply_profile(const char* nic_id, const char* mode, const char* ipv4, const char* gateway, int simulate_failure, char* out, int out_len);
int nas_vm_create(const char* name, int cpu, int memory_gb, int disk_gb, const char* network_id, char* out, int out_len);
int nas_vm_power(const char* vm_id, const char* action, char* out, int out_len);

#endif

