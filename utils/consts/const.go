package consts

const MONGO_DATABASE = "mongo_test"

// jwt related
const JWT_ISSUER = "gin-auth-mongo"
const JWT_ACCESS_TOKEN_EXPIRY = 60 * 24 * 14 // unit: minutes
const JWT_REFRESH_TOKEN_EXPIRY = 90          // unit: days

// email register and reset password
const VERIFY_EMAIL_REGISTER_FLOW_ID = "verify:email:register:flow_id:"
const VERIFY_EMAIL_REGISTER_USERNAME = "verify:email:register:username:"
const VERIFY_EMAIL_REGISTER_PASSWORD = "verify:email:register:password:"
const VERIFY_EMAIL_REGISTER_CODE = "verify:email:register:code:"
const VERIFY_EMAIL_REGISTER_LINK_EXPIRY = 120 // unit: minutes
const VERIFY_EMAIL_REGISTER_CODE_EXPIRY = 15  // unit: minutes

const VERIFY_EMAIL_RESET_PWD_FLOW_ID = "verify:email:resetpwd:flow_id:"
const VERIFY_EMAIL_RESET_PWD_CODE = "verify:email:resetpwd:code:"
const VERIFY_EMAIL_RESET_PWD_PASSWORD = "verify:email:resetpwd:password:"
const VERIFY_EMAIL_RESET_PWD_LINK_EXPIRY = 120 // unit: minutes
const VERIFY_EMAIL_RESET_PWD_CODE_EXPIRY = 15  // unit: minutes

// date and time format
const DATE_FORMAT = "2006-01-02"
const DATETIME_FORMAT = "2006-01-02 15:04:05"
const DATETIME_NANO_FORMAT = "2006-01-02T15:04:05.000" // prefer to use this format
const DATE_FORMAT_REGEX_PATTERN = "^\\d{4}-\\d{2}-\\d{2}$"

// flow limit
const FLOW_LIMIT_PERIOD = 3   // unit: seconds
const FLOW_LIMIT_MAX = 30     // max requests per period
const FLOW_LIMIT_BLOCKED = 30 // unit: minutes // blocked after being limited
const FLOW_LIMIT_COUNTER_KEY = "flow:counter:"
const FLOW_LIMIT_BLOCKED_KEY = "flow:blocked:"

// file upload
const MAX_FILE_SIZE = 500 * 1024 * 1024       // 500MB
const MAX_IMAGE_FILE_SIZE = 250 * 1024 * 1024 // 250MB

// jwk related
const PRIVATE_KEYS_FILE = ".private/keys.json"
const PUBLIC_KEYS_FILE = ".public/keys.json"

// user related
const DEFAULT_AVATAR = MINIO_PUBLIC_BUCKET_NAME + "/avatars/default.svg"
const DEFAULT_COVER_IMAGE = MINIO_PUBLIC_BUCKET_NAME + "/cover_images/default.svg"

const FRONTEND_REGISTER_ROUTE = "/auth/sign-up/complete"
const FRONTEND_RESET_PASSWORD_ROUTE = "/auth/reset-password/complete"

var TRIP_PLAN_USER_PERMISSION_TYPE = []string{"view", "edit", "admin"}

var PARTICIPANT_PERMISSION_VIEW = []string{"view", "edit", "admin"}
var PARTICIPANT_PERMISSION_EDIT = []string{"edit", "admin"}
var PARTICIPANT_PERMISSION_ADMIN = []string{"admin"}
var PARTICIPANT_PERMISSION_OWNER = []string{"owner"}

const TRIP_PLAN_STATUS_UPCOMING = "upcoming"
const TRIP_PLAN_STATUS_ONGOING = "ongoing"
const TRIP_PLAN_STATUS_BRAINSTORMING = "brainstorming"
const TRIP_PLAN_STATUS_FINISHED = "finished"

var TRIP_PLAN_STATUS = []string{TRIP_PLAN_STATUS_UPCOMING, TRIP_PLAN_STATUS_ONGOING, TRIP_PLAN_STATUS_BRAINSTORMING, TRIP_PLAN_STATUS_FINISHED}

const MINIO_PUBLIC_BUCKET_NAME = "mytrip-public"
const MINIO_PRIVATE_BUCKET_NAME = "mytrip-private"

// openssl rand -base64 32 and get the top 32 characters, MUST BE 32 characters
const FILE_ENCRYPTION_KEY = "rxRkW7v2MmO7c0L2jhGTAqN12g+a3fmh"
