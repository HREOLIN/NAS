#include "native.h"

#include <stdio.h>
#include <string.h>

static void write_message(char* out, int out_len, const char* message) {
    if (out == NULL || out_len <= 0) {
        return;
    }
    snprintf(out, out_len, "%s", message);
    out[out_len - 1] = '\0';
}

int nas_recovery_create_snapshot(const char* volume, const char* name, char* out, int out_len) {
    char buffer[256];
    snprintf(buffer, sizeof(buffer), "snapshot created for volume=%s name=%s", volume, name);
    write_message(out, out_len, buffer);
    return 0;
}

int nas_recovery_restore_file(const char* snapshot_id, const char* source_path, const char* target_path, char* out, int out_len) {
    char buffer[256];
    snprintf(buffer, sizeof(buffer), "restore file from snapshot=%s source=%s target=%s", snapshot_id, source_path, target_path);
    write_message(out, out_len, buffer);
    return 0;
}

int nas_recovery_rollback_volume(const char* snapshot_id, const char* volume, char* out, int out_len) {
    char buffer[256];
    snprintf(buffer, sizeof(buffer), "rollback volume=%s using snapshot=%s", volume, snapshot_id);
    write_message(out, out_len, buffer);
    return 0;
}

int nas_network_apply_profile(const char* platform, const char* nic_id, const char* commands_blob, const char* probes_blob, int simulate_failure, char* out, int out_len) {
    if (simulate_failure != 0) {
        write_message(out, out_len, "connectivity probe failed, rollback required");
        return 2;
    }

    char buffer[256];
    snprintf(buffer, sizeof(buffer), "network plan applied platform=%s nic=%s commands=%s probes=%s", platform, nic_id, commands_blob, probes_blob);
    write_message(out, out_len, buffer);
    return 0;
}

int nas_vm_create(const char* name, int cpu, int memory_gb, int disk_gb, const char* network_id, char* out, int out_len) {
    char buffer[256];
    snprintf(buffer, sizeof(buffer), "vm prepared name=%s cpu=%d memory=%dGB disk=%dGB network=%s", name, cpu, memory_gb, disk_gb, network_id);
    write_message(out, out_len, buffer);
    return 0;
}

int nas_vm_power(const char* vm_id, const char* action, char* out, int out_len) {
    if (strcmp(action, "start") != 0 && strcmp(action, "stop") != 0 && strcmp(action, "restart") != 0) {
        write_message(out, out_len, "unsupported vm power action");
        return 3;
    }

    char buffer[256];
    snprintf(buffer, sizeof(buffer), "vm power action applied vm=%s action=%s", vm_id, action);
    write_message(out, out_len, buffer);
    return 0;
}
