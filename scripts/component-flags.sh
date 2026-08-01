#!/usr/bin/env bash

component_disabled() {
  local root="$1"
  local component="$2"
  [ -f "$root/db/${component}.disabled" ]
}
