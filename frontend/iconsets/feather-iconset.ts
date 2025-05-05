import { h } from "vue"
import type { IconSet, IconAliases, IconProps } from "vuetify"

const aliases: IconAliases = {
  calendar: "fr-calendar", // calendar: mdi-calendar
  cancel: "fr-x-circle", // cancel: mdi-close-circle
  checkboxIndeterminate: "fr-minus-square", // checkboxIndeterminate: mdi-minus-box
  checkboxOff: "fr-square", // checkboxOff: mdi-checkbox-blank-outline
  checkboxOn: "fr-check-square-filled", // checkboxOn: mdi-checkbox-marked
  clear: "fr-x-circle", // clear: mdi-close-circle
  close: "fr-x", // close: mdi-close
  collapse: "fr-chevron-up", // collapse: mdi-chevron-up
  complete: "fr-check", // complete: mdi-check
  delete: "fr-x-circle", // delete: mdi-close-circle
  delimiter: "fr-circle", // delimiter: mdi-circle
  dropdown: "fr-chevron-down", // dropdown: mdi-menu-down
  edit: "fr-edit", // edit: mdi-pencil
  error: "fr-x-circle", // error: mdi-close-circle
  expand: "fr-chevron-down", // expand: mdi-chevron-down
  eyeDropper: "fr-pen-tool", // eyeDropper: mdi-eyedropper
  file: "fr-paperclip", // file: mdi-paperclip
  first: "fr-chevrons-left", // first: mdi-page-first
  info: "fr-info", // info: mdi-information
  last: "fr-chevrons-right", // last: mdi-page-last
  loading: "fr-refresh-cw", // loading: mdi-cached
  menu: "fr-menu", // menu: mdi-menu
  minus: "fr-minus", // minus: mdi-minus
  next: "fr-chevron-right", // next: mdi-chevron-right
  plus: "fr-plus", // plus: mdi-plus
  prev: "fr-chevron-left", // prev: mdi-chevron-left
  radioOff: "fr-circle", // radioOff: mdi-radiobox-blank
  radioOn: "fr-check-circle", // radioOn: mdi-radiobox-marked
  ratingEmpty: "fr-star", // ratingEmpty: mdi-star-outline
  ratingFull: "fr-star", // ratingFull: mdi-star
  ratingHalf: "fr-star", // ratingHalf: mdi-star-half-full
  sortAsc: "fr-arrow-up", // sortAsc: mdi-arrow-up
  sortDesc: "fr-arrow-down", // sortDesc: mdi-arrow-down
  subgroup: "fr-corner-right-down", // subgroup: mdi-menu-down
  success: "fr-check-circle", // success: mdi-check-circle
  unfold: "fr-corner-right-down", // unfold: mdi-unfold-more-horizontal
  warning: "fr-alert-circle", // warning: mdi-alert-circle
}

const fr: IconSet = {
  component: (props: IconProps) =>
    h("i", {
      ...props,
      class: props.icon,
    }),
}

export { aliases, fr }
