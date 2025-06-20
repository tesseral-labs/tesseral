// @generated by protoc-gen-es v2.2.3 with parameter "target=ts"
// @generated from file tesseral/backend/v1/audit_log_events.proto (package tesseral.backend.v1, syntax proto3)
/* eslint-disable */

import type { GenFile, GenMessage } from "@bufbuild/protobuf/codegenv1";
import { fileDesc, messageDesc } from "@bufbuild/protobuf/codegenv1";
import type { APIKey, APIKeyRoleAssignment, Organization, Passkey, Role, SAMLConnection, SCIMAPIKey, User, UserInvite, UserRoleAssignment } from "./models_pb";
import { file_tesseral_backend_v1_models } from "./models_pb";
import type { Message } from "@bufbuild/protobuf";

/**
 * Describes the file tesseral/backend/v1/audit_log_events.proto.
 */
export const file_tesseral_backend_v1_audit_log_events: GenFile = /*@__PURE__*/
  fileDesc("Cip0ZXNzZXJhbC9iYWNrZW5kL3YxL2F1ZGl0X2xvZ19ldmVudHMucHJvdG8SE3Rlc3NlcmFsLmJhY2tlbmQudjEiYQobQVBJS2V5Um9sZUFzc2lnbm1lbnRDcmVhdGVkEkIKD3JvbGVfYXNzaWdubWVudBgBIAEoCzIpLnRlc3NlcmFsLmJhY2tlbmQudjEuQVBJS2V5Um9sZUFzc2lnbm1lbnQiYQobQVBJS2V5Um9sZUFzc2lnbm1lbnREZWxldGVkEkIKD3JvbGVfYXNzaWdubWVudBgBIAEoCzIpLnRlc3NlcmFsLmJhY2tlbmQudjEuQVBJS2V5Um9sZUFzc2lnbm1lbnQirgEKG0FQSUtleVJvbGVBc3NpZ25tZW50VXBkYXRlZBJCCg9yb2xlX2Fzc2lnbm1lbnQYASABKAsyKS50ZXNzZXJhbC5iYWNrZW5kLnYxLkFQSUtleVJvbGVBc3NpZ25tZW50EksKGHByZXZpb3VzX3JvbGVfYXNzaWdubWVudBgCIAEoCzIpLnRlc3NlcmFsLmJhY2tlbmQudjEuQVBJS2V5Um9sZUFzc2lnbm1lbnQiPQoNQVBJS2V5Q3JlYXRlZBIsCgdhcGlfa2V5GAEgASgLMhsudGVzc2VyYWwuYmFja2VuZC52MS5BUElLZXkiPQoNQVBJS2V5RGVsZXRlZBIsCgdhcGlfa2V5GAEgASgLMhsudGVzc2VyYWwuYmFja2VuZC52MS5BUElLZXkidAoNQVBJS2V5UmV2b2tlZBIsCgdhcGlfa2V5GAEgASgLMhsudGVzc2VyYWwuYmFja2VuZC52MS5BUElLZXkSNQoQcHJldmlvdXNfYXBpX2tleRgCIAEoCzIbLnRlc3NlcmFsLmJhY2tlbmQudjEuQVBJS2V5InQKDUFQSUtleVVwZGF0ZWQSLAoHYXBpX2tleRgBIAEoCzIbLnRlc3NlcmFsLmJhY2tlbmQudjEuQVBJS2V5EjUKEHByZXZpb3VzX2FwaV9rZXkYAiABKAsyGy50ZXNzZXJhbC5iYWNrZW5kLnYxLkFQSUtleSJOChNPcmdhbml6YXRpb25DcmVhdGVkEjcKDG9yZ2FuaXphdGlvbhgBIAEoCzIhLnRlc3NlcmFsLmJhY2tlbmQudjEuT3JnYW5pemF0aW9uIk4KE09yZ2FuaXphdGlvbkRlbGV0ZWQSNwoMb3JnYW5pemF0aW9uGAEgASgLMiEudGVzc2VyYWwuYmFja2VuZC52MS5Pcmdhbml6YXRpb24ikAEKE09yZ2FuaXphdGlvblVwZGF0ZWQSNwoMb3JnYW5pemF0aW9uGAEgASgLMiEudGVzc2VyYWwuYmFja2VuZC52MS5Pcmdhbml6YXRpb24SQAoVcHJldmlvdXNfb3JnYW5pemF0aW9uGAIgASgLMiEudGVzc2VyYWwuYmFja2VuZC52MS5Pcmdhbml6YXRpb24ibwomT3JnYW5pemF0aW9uR29vZ2xlSG9zdGVkRG9tYWluc1VwZGF0ZWQSHQoVZ29vZ2xlX2hvc3RlZF9kb21haW5zGAEgAygJEiYKHnByZXZpb3VzX2dvb2dsZV9ob3N0ZWRfZG9tYWlucxgCIAMoCSJHChpPcmdhbml6YXRpb25Eb21haW5zVXBkYXRlZBIPCgdkb21haW5zGAEgAygJEhgKEHByZXZpb3VzX2RvbWFpbnMYAiADKAkibAolT3JnYW5pemF0aW9uTWljcm9zb2Z0VGVuYW50SURzVXBkYXRlZBIcChRtaWNyb3NvZnRfdGVuYW50X2lkcxgBIAMoCRIlCh1wcmV2aW91c19taWNyb3NvZnRfdGVuYW50X2lkcxgCIAMoCSI/Cg5QYXNza2V5Q3JlYXRlZBItCgdwYXNza2V5GAEgASgLMhwudGVzc2VyYWwuYmFja2VuZC52MS5QYXNza2V5Ij8KDlBhc3NrZXlEZWxldGVkEi0KB3Bhc3NrZXkYASABKAsyHC50ZXNzZXJhbC5iYWNrZW5kLnYxLlBhc3NrZXkidwoOUGFzc2tleVVwZGF0ZWQSLQoHcGFzc2tleRgBIAEoCzIcLnRlc3NlcmFsLmJhY2tlbmQudjEuUGFzc2tleRI2ChBwcmV2aW91c19wYXNza2V5GAIgASgLMhwudGVzc2VyYWwuYmFja2VuZC52MS5QYXNza2V5IjYKC1JvbGVDcmVhdGVkEicKBHJvbGUYASABKAsyGS50ZXNzZXJhbC5iYWNrZW5kLnYxLlJvbGUiNgoLUm9sZURlbGV0ZWQSJwoEcm9sZRgBIAEoCzIZLnRlc3NlcmFsLmJhY2tlbmQudjEuUm9sZSJoCgtSb2xlVXBkYXRlZBInCgRyb2xlGAEgASgLMhkudGVzc2VyYWwuYmFja2VuZC52MS5Sb2xlEjAKDXByZXZpb3VzX3JvbGUYAiABKAsyGS50ZXNzZXJhbC5iYWNrZW5kLnYxLlJvbGUiVQoVU0FNTENvbm5lY3Rpb25DcmVhdGVkEjwKD3NhbWxfY29ubmVjdGlvbhgBIAEoCzIjLnRlc3NlcmFsLmJhY2tlbmQudjEuU0FNTENvbm5lY3Rpb24iVQoVU0FNTENvbm5lY3Rpb25EZWxldGVkEjwKD3NhbWxfY29ubmVjdGlvbhgBIAEoCzIjLnRlc3NlcmFsLmJhY2tlbmQudjEuU0FNTENvbm5lY3Rpb24inAEKFVNBTUxDb25uZWN0aW9uVXBkYXRlZBI8Cg9zYW1sX2Nvbm5lY3Rpb24YASABKAsyIy50ZXNzZXJhbC5iYWNrZW5kLnYxLlNBTUxDb25uZWN0aW9uEkUKGHByZXZpb3VzX3NhbWxfY29ubmVjdGlvbhgCIAEoCzIjLnRlc3NlcmFsLmJhY2tlbmQudjEuU0FNTENvbm5lY3Rpb24iSgoRU0NJTUFQSUtleUNyZWF0ZWQSNQoMc2NpbV9hcGlfa2V5GAEgASgLMh8udGVzc2VyYWwuYmFja2VuZC52MS5TQ0lNQVBJS2V5IkoKEVNDSU1BUElLZXlEZWxldGVkEjUKDHNjaW1fYXBpX2tleRgBIAEoCzIfLnRlc3NlcmFsLmJhY2tlbmQudjEuU0NJTUFQSUtleSKKAQoRU0NJTUFQSUtleVJldm9rZWQSNQoMc2NpbV9hcGlfa2V5GAEgASgLMh8udGVzc2VyYWwuYmFja2VuZC52MS5TQ0lNQVBJS2V5Ej4KFXByZXZpb3VzX3NjaW1fYXBpX2tleRgCIAEoCzIfLnRlc3NlcmFsLmJhY2tlbmQudjEuU0NJTUFQSUtleSKKAQoRU0NJTUFQSUtleVVwZGF0ZWQSNQoMc2NpbV9hcGlfa2V5GAEgASgLMh8udGVzc2VyYWwuYmFja2VuZC52MS5TQ0lNQVBJS2V5Ej4KFXByZXZpb3VzX3NjaW1fYXBpX2tleRgCIAEoCzIfLnRlc3NlcmFsLmJhY2tlbmQudjEuU0NJTUFQSUtleSJJChFVc2VySW52aXRlQ3JlYXRlZBI0Cgt1c2VyX2ludml0ZRgBIAEoCzIfLnRlc3NlcmFsLmJhY2tlbmQudjEuVXNlckludml0ZSJJChFVc2VySW52aXRlRGVsZXRlZBI0Cgt1c2VyX2ludml0ZRgBIAEoCzIfLnRlc3NlcmFsLmJhY2tlbmQudjEuVXNlckludml0ZSKIAQoRVXNlckludml0ZVVwZGF0ZWQSNAoLdXNlcl9pbnZpdGUYASABKAsyHy50ZXNzZXJhbC5iYWNrZW5kLnYxLlVzZXJJbnZpdGUSPQoUcHJldmlvdXNfdXNlcl9pbnZpdGUYAiABKAsyHy50ZXNzZXJhbC5iYWNrZW5kLnYxLlVzZXJJbnZpdGUiYgoZVXNlclJvbGVBc3NpZ25tZW50Q3JlYXRlZBJFChR1c2VyX3JvbGVfYXNzaWdubWVudBgBIAEoCzInLnRlc3NlcmFsLmJhY2tlbmQudjEuVXNlclJvbGVBc3NpZ25tZW50ImIKGVVzZXJSb2xlQXNzaWdubWVudERlbGV0ZWQSRQoUdXNlcl9yb2xlX2Fzc2lnbm1lbnQYASABKAsyJy50ZXNzZXJhbC5iYWNrZW5kLnYxLlVzZXJSb2xlQXNzaWdubWVudCI2CgtVc2VyQ3JlYXRlZBInCgR1c2VyGAEgASgLMhkudGVzc2VyYWwuYmFja2VuZC52MS5Vc2VyIjYKC1VzZXJEZWxldGVkEicKBHVzZXIYASABKAsyGS50ZXNzZXJhbC5iYWNrZW5kLnYxLlVzZXIiaAoLVXNlclVwZGF0ZWQSJwoEdXNlchgBIAEoCzIZLnRlc3NlcmFsLmJhY2tlbmQudjEuVXNlchIwCg1wcmV2aW91c191c2VyGAIgASgLMhkudGVzc2VyYWwuYmFja2VuZC52MS5Vc2VyQvIBChdjb20udGVzc2VyYWwuYmFja2VuZC52MUITQXVkaXRMb2dFdmVudHNQcm90b1ABWlRnaXRodWIuY29tL3Rlc3NlcmFsLWxhYnMvdGVzc2VyYWwvaW50ZXJuYWwvYmFja2VuZC9nZW4vdGVzc2VyYWwvYmFja2VuZC92MTtiYWNrZW5kdjGiAgNUQliqAhNUZXNzZXJhbC5CYWNrZW5kLlYxygITVGVzc2VyYWxcQmFja2VuZFxWMeICH1Rlc3NlcmFsXEJhY2tlbmRcVjFcR1BCTWV0YWRhdGHqAhVUZXNzZXJhbDo6QmFja2VuZDo6VjFiBnByb3RvMw", [file_tesseral_backend_v1_models]);

/**
 * @generated from message tesseral.backend.v1.APIKeyRoleAssignmentCreated
 */
export type APIKeyRoleAssignmentCreated = Message<"tesseral.backend.v1.APIKeyRoleAssignmentCreated"> & {
  /**
   * @generated from field: tesseral.backend.v1.APIKeyRoleAssignment role_assignment = 1;
   */
  roleAssignment?: APIKeyRoleAssignment;
};

/**
 * Describes the message tesseral.backend.v1.APIKeyRoleAssignmentCreated.
 * Use `create(APIKeyRoleAssignmentCreatedSchema)` to create a new message.
 */
export const APIKeyRoleAssignmentCreatedSchema: GenMessage<APIKeyRoleAssignmentCreated> = /*@__PURE__*/
  messageDesc(file_tesseral_backend_v1_audit_log_events, 0);

/**
 * @generated from message tesseral.backend.v1.APIKeyRoleAssignmentDeleted
 */
export type APIKeyRoleAssignmentDeleted = Message<"tesseral.backend.v1.APIKeyRoleAssignmentDeleted"> & {
  /**
   * @generated from field: tesseral.backend.v1.APIKeyRoleAssignment role_assignment = 1;
   */
  roleAssignment?: APIKeyRoleAssignment;
};

/**
 * Describes the message tesseral.backend.v1.APIKeyRoleAssignmentDeleted.
 * Use `create(APIKeyRoleAssignmentDeletedSchema)` to create a new message.
 */
export const APIKeyRoleAssignmentDeletedSchema: GenMessage<APIKeyRoleAssignmentDeleted> = /*@__PURE__*/
  messageDesc(file_tesseral_backend_v1_audit_log_events, 1);

/**
 * @generated from message tesseral.backend.v1.APIKeyRoleAssignmentUpdated
 */
export type APIKeyRoleAssignmentUpdated = Message<"tesseral.backend.v1.APIKeyRoleAssignmentUpdated"> & {
  /**
   * @generated from field: tesseral.backend.v1.APIKeyRoleAssignment role_assignment = 1;
   */
  roleAssignment?: APIKeyRoleAssignment;

  /**
   * @generated from field: tesseral.backend.v1.APIKeyRoleAssignment previous_role_assignment = 2;
   */
  previousRoleAssignment?: APIKeyRoleAssignment;
};

/**
 * Describes the message tesseral.backend.v1.APIKeyRoleAssignmentUpdated.
 * Use `create(APIKeyRoleAssignmentUpdatedSchema)` to create a new message.
 */
export const APIKeyRoleAssignmentUpdatedSchema: GenMessage<APIKeyRoleAssignmentUpdated> = /*@__PURE__*/
  messageDesc(file_tesseral_backend_v1_audit_log_events, 2);

/**
 * @generated from message tesseral.backend.v1.APIKeyCreated
 */
export type APIKeyCreated = Message<"tesseral.backend.v1.APIKeyCreated"> & {
  /**
   * @generated from field: tesseral.backend.v1.APIKey api_key = 1;
   */
  apiKey?: APIKey;
};

/**
 * Describes the message tesseral.backend.v1.APIKeyCreated.
 * Use `create(APIKeyCreatedSchema)` to create a new message.
 */
export const APIKeyCreatedSchema: GenMessage<APIKeyCreated> = /*@__PURE__*/
  messageDesc(file_tesseral_backend_v1_audit_log_events, 3);

/**
 * @generated from message tesseral.backend.v1.APIKeyDeleted
 */
export type APIKeyDeleted = Message<"tesseral.backend.v1.APIKeyDeleted"> & {
  /**
   * @generated from field: tesseral.backend.v1.APIKey api_key = 1;
   */
  apiKey?: APIKey;
};

/**
 * Describes the message tesseral.backend.v1.APIKeyDeleted.
 * Use `create(APIKeyDeletedSchema)` to create a new message.
 */
export const APIKeyDeletedSchema: GenMessage<APIKeyDeleted> = /*@__PURE__*/
  messageDesc(file_tesseral_backend_v1_audit_log_events, 4);

/**
 * @generated from message tesseral.backend.v1.APIKeyRevoked
 */
export type APIKeyRevoked = Message<"tesseral.backend.v1.APIKeyRevoked"> & {
  /**
   * @generated from field: tesseral.backend.v1.APIKey api_key = 1;
   */
  apiKey?: APIKey;

  /**
   * @generated from field: tesseral.backend.v1.APIKey previous_api_key = 2;
   */
  previousApiKey?: APIKey;
};

/**
 * Describes the message tesseral.backend.v1.APIKeyRevoked.
 * Use `create(APIKeyRevokedSchema)` to create a new message.
 */
export const APIKeyRevokedSchema: GenMessage<APIKeyRevoked> = /*@__PURE__*/
  messageDesc(file_tesseral_backend_v1_audit_log_events, 5);

/**
 * @generated from message tesseral.backend.v1.APIKeyUpdated
 */
export type APIKeyUpdated = Message<"tesseral.backend.v1.APIKeyUpdated"> & {
  /**
   * @generated from field: tesseral.backend.v1.APIKey api_key = 1;
   */
  apiKey?: APIKey;

  /**
   * @generated from field: tesseral.backend.v1.APIKey previous_api_key = 2;
   */
  previousApiKey?: APIKey;
};

/**
 * Describes the message tesseral.backend.v1.APIKeyUpdated.
 * Use `create(APIKeyUpdatedSchema)` to create a new message.
 */
export const APIKeyUpdatedSchema: GenMessage<APIKeyUpdated> = /*@__PURE__*/
  messageDesc(file_tesseral_backend_v1_audit_log_events, 6);

/**
 * @generated from message tesseral.backend.v1.OrganizationCreated
 */
export type OrganizationCreated = Message<"tesseral.backend.v1.OrganizationCreated"> & {
  /**
   * @generated from field: tesseral.backend.v1.Organization organization = 1;
   */
  organization?: Organization;
};

/**
 * Describes the message tesseral.backend.v1.OrganizationCreated.
 * Use `create(OrganizationCreatedSchema)` to create a new message.
 */
export const OrganizationCreatedSchema: GenMessage<OrganizationCreated> = /*@__PURE__*/
  messageDesc(file_tesseral_backend_v1_audit_log_events, 7);

/**
 * @generated from message tesseral.backend.v1.OrganizationDeleted
 */
export type OrganizationDeleted = Message<"tesseral.backend.v1.OrganizationDeleted"> & {
  /**
   * @generated from field: tesseral.backend.v1.Organization organization = 1;
   */
  organization?: Organization;
};

/**
 * Describes the message tesseral.backend.v1.OrganizationDeleted.
 * Use `create(OrganizationDeletedSchema)` to create a new message.
 */
export const OrganizationDeletedSchema: GenMessage<OrganizationDeleted> = /*@__PURE__*/
  messageDesc(file_tesseral_backend_v1_audit_log_events, 8);

/**
 * @generated from message tesseral.backend.v1.OrganizationUpdated
 */
export type OrganizationUpdated = Message<"tesseral.backend.v1.OrganizationUpdated"> & {
  /**
   * @generated from field: tesseral.backend.v1.Organization organization = 1;
   */
  organization?: Organization;

  /**
   * @generated from field: tesseral.backend.v1.Organization previous_organization = 2;
   */
  previousOrganization?: Organization;
};

/**
 * Describes the message tesseral.backend.v1.OrganizationUpdated.
 * Use `create(OrganizationUpdatedSchema)` to create a new message.
 */
export const OrganizationUpdatedSchema: GenMessage<OrganizationUpdated> = /*@__PURE__*/
  messageDesc(file_tesseral_backend_v1_audit_log_events, 9);

/**
 * @generated from message tesseral.backend.v1.OrganizationGoogleHostedDomainsUpdated
 */
export type OrganizationGoogleHostedDomainsUpdated = Message<"tesseral.backend.v1.OrganizationGoogleHostedDomainsUpdated"> & {
  /**
   * @generated from field: repeated string google_hosted_domains = 1;
   */
  googleHostedDomains: string[];

  /**
   * @generated from field: repeated string previous_google_hosted_domains = 2;
   */
  previousGoogleHostedDomains: string[];
};

/**
 * Describes the message tesseral.backend.v1.OrganizationGoogleHostedDomainsUpdated.
 * Use `create(OrganizationGoogleHostedDomainsUpdatedSchema)` to create a new message.
 */
export const OrganizationGoogleHostedDomainsUpdatedSchema: GenMessage<OrganizationGoogleHostedDomainsUpdated> = /*@__PURE__*/
  messageDesc(file_tesseral_backend_v1_audit_log_events, 10);

/**
 * @generated from message tesseral.backend.v1.OrganizationDomainsUpdated
 */
export type OrganizationDomainsUpdated = Message<"tesseral.backend.v1.OrganizationDomainsUpdated"> & {
  /**
   * @generated from field: repeated string domains = 1;
   */
  domains: string[];

  /**
   * @generated from field: repeated string previous_domains = 2;
   */
  previousDomains: string[];
};

/**
 * Describes the message tesseral.backend.v1.OrganizationDomainsUpdated.
 * Use `create(OrganizationDomainsUpdatedSchema)` to create a new message.
 */
export const OrganizationDomainsUpdatedSchema: GenMessage<OrganizationDomainsUpdated> = /*@__PURE__*/
  messageDesc(file_tesseral_backend_v1_audit_log_events, 11);

/**
 * @generated from message tesseral.backend.v1.OrganizationMicrosoftTenantIDsUpdated
 */
export type OrganizationMicrosoftTenantIDsUpdated = Message<"tesseral.backend.v1.OrganizationMicrosoftTenantIDsUpdated"> & {
  /**
   * @generated from field: repeated string microsoft_tenant_ids = 1;
   */
  microsoftTenantIds: string[];

  /**
   * @generated from field: repeated string previous_microsoft_tenant_ids = 2;
   */
  previousMicrosoftTenantIds: string[];
};

/**
 * Describes the message tesseral.backend.v1.OrganizationMicrosoftTenantIDsUpdated.
 * Use `create(OrganizationMicrosoftTenantIDsUpdatedSchema)` to create a new message.
 */
export const OrganizationMicrosoftTenantIDsUpdatedSchema: GenMessage<OrganizationMicrosoftTenantIDsUpdated> = /*@__PURE__*/
  messageDesc(file_tesseral_backend_v1_audit_log_events, 12);

/**
 * @generated from message tesseral.backend.v1.PasskeyCreated
 */
export type PasskeyCreated = Message<"tesseral.backend.v1.PasskeyCreated"> & {
  /**
   * @generated from field: tesseral.backend.v1.Passkey passkey = 1;
   */
  passkey?: Passkey;
};

/**
 * Describes the message tesseral.backend.v1.PasskeyCreated.
 * Use `create(PasskeyCreatedSchema)` to create a new message.
 */
export const PasskeyCreatedSchema: GenMessage<PasskeyCreated> = /*@__PURE__*/
  messageDesc(file_tesseral_backend_v1_audit_log_events, 13);

/**
 * @generated from message tesseral.backend.v1.PasskeyDeleted
 */
export type PasskeyDeleted = Message<"tesseral.backend.v1.PasskeyDeleted"> & {
  /**
   * @generated from field: tesseral.backend.v1.Passkey passkey = 1;
   */
  passkey?: Passkey;
};

/**
 * Describes the message tesseral.backend.v1.PasskeyDeleted.
 * Use `create(PasskeyDeletedSchema)` to create a new message.
 */
export const PasskeyDeletedSchema: GenMessage<PasskeyDeleted> = /*@__PURE__*/
  messageDesc(file_tesseral_backend_v1_audit_log_events, 14);

/**
 * @generated from message tesseral.backend.v1.PasskeyUpdated
 */
export type PasskeyUpdated = Message<"tesseral.backend.v1.PasskeyUpdated"> & {
  /**
   * @generated from field: tesseral.backend.v1.Passkey passkey = 1;
   */
  passkey?: Passkey;

  /**
   * @generated from field: tesseral.backend.v1.Passkey previous_passkey = 2;
   */
  previousPasskey?: Passkey;
};

/**
 * Describes the message tesseral.backend.v1.PasskeyUpdated.
 * Use `create(PasskeyUpdatedSchema)` to create a new message.
 */
export const PasskeyUpdatedSchema: GenMessage<PasskeyUpdated> = /*@__PURE__*/
  messageDesc(file_tesseral_backend_v1_audit_log_events, 15);

/**
 * @generated from message tesseral.backend.v1.RoleCreated
 */
export type RoleCreated = Message<"tesseral.backend.v1.RoleCreated"> & {
  /**
   * @generated from field: tesseral.backend.v1.Role role = 1;
   */
  role?: Role;
};

/**
 * Describes the message tesseral.backend.v1.RoleCreated.
 * Use `create(RoleCreatedSchema)` to create a new message.
 */
export const RoleCreatedSchema: GenMessage<RoleCreated> = /*@__PURE__*/
  messageDesc(file_tesseral_backend_v1_audit_log_events, 16);

/**
 * @generated from message tesseral.backend.v1.RoleDeleted
 */
export type RoleDeleted = Message<"tesseral.backend.v1.RoleDeleted"> & {
  /**
   * @generated from field: tesseral.backend.v1.Role role = 1;
   */
  role?: Role;
};

/**
 * Describes the message tesseral.backend.v1.RoleDeleted.
 * Use `create(RoleDeletedSchema)` to create a new message.
 */
export const RoleDeletedSchema: GenMessage<RoleDeleted> = /*@__PURE__*/
  messageDesc(file_tesseral_backend_v1_audit_log_events, 17);

/**
 * @generated from message tesseral.backend.v1.RoleUpdated
 */
export type RoleUpdated = Message<"tesseral.backend.v1.RoleUpdated"> & {
  /**
   * @generated from field: tesseral.backend.v1.Role role = 1;
   */
  role?: Role;

  /**
   * @generated from field: tesseral.backend.v1.Role previous_role = 2;
   */
  previousRole?: Role;
};

/**
 * Describes the message tesseral.backend.v1.RoleUpdated.
 * Use `create(RoleUpdatedSchema)` to create a new message.
 */
export const RoleUpdatedSchema: GenMessage<RoleUpdated> = /*@__PURE__*/
  messageDesc(file_tesseral_backend_v1_audit_log_events, 18);

/**
 * @generated from message tesseral.backend.v1.SAMLConnectionCreated
 */
export type SAMLConnectionCreated = Message<"tesseral.backend.v1.SAMLConnectionCreated"> & {
  /**
   * @generated from field: tesseral.backend.v1.SAMLConnection saml_connection = 1;
   */
  samlConnection?: SAMLConnection;
};

/**
 * Describes the message tesseral.backend.v1.SAMLConnectionCreated.
 * Use `create(SAMLConnectionCreatedSchema)` to create a new message.
 */
export const SAMLConnectionCreatedSchema: GenMessage<SAMLConnectionCreated> = /*@__PURE__*/
  messageDesc(file_tesseral_backend_v1_audit_log_events, 19);

/**
 * @generated from message tesseral.backend.v1.SAMLConnectionDeleted
 */
export type SAMLConnectionDeleted = Message<"tesseral.backend.v1.SAMLConnectionDeleted"> & {
  /**
   * @generated from field: tesseral.backend.v1.SAMLConnection saml_connection = 1;
   */
  samlConnection?: SAMLConnection;
};

/**
 * Describes the message tesseral.backend.v1.SAMLConnectionDeleted.
 * Use `create(SAMLConnectionDeletedSchema)` to create a new message.
 */
export const SAMLConnectionDeletedSchema: GenMessage<SAMLConnectionDeleted> = /*@__PURE__*/
  messageDesc(file_tesseral_backend_v1_audit_log_events, 20);

/**
 * @generated from message tesseral.backend.v1.SAMLConnectionUpdated
 */
export type SAMLConnectionUpdated = Message<"tesseral.backend.v1.SAMLConnectionUpdated"> & {
  /**
   * @generated from field: tesseral.backend.v1.SAMLConnection saml_connection = 1;
   */
  samlConnection?: SAMLConnection;

  /**
   * @generated from field: tesseral.backend.v1.SAMLConnection previous_saml_connection = 2;
   */
  previousSamlConnection?: SAMLConnection;
};

/**
 * Describes the message tesseral.backend.v1.SAMLConnectionUpdated.
 * Use `create(SAMLConnectionUpdatedSchema)` to create a new message.
 */
export const SAMLConnectionUpdatedSchema: GenMessage<SAMLConnectionUpdated> = /*@__PURE__*/
  messageDesc(file_tesseral_backend_v1_audit_log_events, 21);

/**
 * @generated from message tesseral.backend.v1.SCIMAPIKeyCreated
 */
export type SCIMAPIKeyCreated = Message<"tesseral.backend.v1.SCIMAPIKeyCreated"> & {
  /**
   * @generated from field: tesseral.backend.v1.SCIMAPIKey scim_api_key = 1;
   */
  scimApiKey?: SCIMAPIKey;
};

/**
 * Describes the message tesseral.backend.v1.SCIMAPIKeyCreated.
 * Use `create(SCIMAPIKeyCreatedSchema)` to create a new message.
 */
export const SCIMAPIKeyCreatedSchema: GenMessage<SCIMAPIKeyCreated> = /*@__PURE__*/
  messageDesc(file_tesseral_backend_v1_audit_log_events, 22);

/**
 * @generated from message tesseral.backend.v1.SCIMAPIKeyDeleted
 */
export type SCIMAPIKeyDeleted = Message<"tesseral.backend.v1.SCIMAPIKeyDeleted"> & {
  /**
   * @generated from field: tesseral.backend.v1.SCIMAPIKey scim_api_key = 1;
   */
  scimApiKey?: SCIMAPIKey;
};

/**
 * Describes the message tesseral.backend.v1.SCIMAPIKeyDeleted.
 * Use `create(SCIMAPIKeyDeletedSchema)` to create a new message.
 */
export const SCIMAPIKeyDeletedSchema: GenMessage<SCIMAPIKeyDeleted> = /*@__PURE__*/
  messageDesc(file_tesseral_backend_v1_audit_log_events, 23);

/**
 * @generated from message tesseral.backend.v1.SCIMAPIKeyRevoked
 */
export type SCIMAPIKeyRevoked = Message<"tesseral.backend.v1.SCIMAPIKeyRevoked"> & {
  /**
   * @generated from field: tesseral.backend.v1.SCIMAPIKey scim_api_key = 1;
   */
  scimApiKey?: SCIMAPIKey;

  /**
   * @generated from field: tesseral.backend.v1.SCIMAPIKey previous_scim_api_key = 2;
   */
  previousScimApiKey?: SCIMAPIKey;
};

/**
 * Describes the message tesseral.backend.v1.SCIMAPIKeyRevoked.
 * Use `create(SCIMAPIKeyRevokedSchema)` to create a new message.
 */
export const SCIMAPIKeyRevokedSchema: GenMessage<SCIMAPIKeyRevoked> = /*@__PURE__*/
  messageDesc(file_tesseral_backend_v1_audit_log_events, 24);

/**
 * @generated from message tesseral.backend.v1.SCIMAPIKeyUpdated
 */
export type SCIMAPIKeyUpdated = Message<"tesseral.backend.v1.SCIMAPIKeyUpdated"> & {
  /**
   * @generated from field: tesseral.backend.v1.SCIMAPIKey scim_api_key = 1;
   */
  scimApiKey?: SCIMAPIKey;

  /**
   * @generated from field: tesseral.backend.v1.SCIMAPIKey previous_scim_api_key = 2;
   */
  previousScimApiKey?: SCIMAPIKey;
};

/**
 * Describes the message tesseral.backend.v1.SCIMAPIKeyUpdated.
 * Use `create(SCIMAPIKeyUpdatedSchema)` to create a new message.
 */
export const SCIMAPIKeyUpdatedSchema: GenMessage<SCIMAPIKeyUpdated> = /*@__PURE__*/
  messageDesc(file_tesseral_backend_v1_audit_log_events, 25);

/**
 * @generated from message tesseral.backend.v1.UserInviteCreated
 */
export type UserInviteCreated = Message<"tesseral.backend.v1.UserInviteCreated"> & {
  /**
   * @generated from field: tesseral.backend.v1.UserInvite user_invite = 1;
   */
  userInvite?: UserInvite;
};

/**
 * Describes the message tesseral.backend.v1.UserInviteCreated.
 * Use `create(UserInviteCreatedSchema)` to create a new message.
 */
export const UserInviteCreatedSchema: GenMessage<UserInviteCreated> = /*@__PURE__*/
  messageDesc(file_tesseral_backend_v1_audit_log_events, 26);

/**
 * @generated from message tesseral.backend.v1.UserInviteDeleted
 */
export type UserInviteDeleted = Message<"tesseral.backend.v1.UserInviteDeleted"> & {
  /**
   * @generated from field: tesseral.backend.v1.UserInvite user_invite = 1;
   */
  userInvite?: UserInvite;
};

/**
 * Describes the message tesseral.backend.v1.UserInviteDeleted.
 * Use `create(UserInviteDeletedSchema)` to create a new message.
 */
export const UserInviteDeletedSchema: GenMessage<UserInviteDeleted> = /*@__PURE__*/
  messageDesc(file_tesseral_backend_v1_audit_log_events, 27);

/**
 * @generated from message tesseral.backend.v1.UserInviteUpdated
 */
export type UserInviteUpdated = Message<"tesseral.backend.v1.UserInviteUpdated"> & {
  /**
   * @generated from field: tesseral.backend.v1.UserInvite user_invite = 1;
   */
  userInvite?: UserInvite;

  /**
   * @generated from field: tesseral.backend.v1.UserInvite previous_user_invite = 2;
   */
  previousUserInvite?: UserInvite;
};

/**
 * Describes the message tesseral.backend.v1.UserInviteUpdated.
 * Use `create(UserInviteUpdatedSchema)` to create a new message.
 */
export const UserInviteUpdatedSchema: GenMessage<UserInviteUpdated> = /*@__PURE__*/
  messageDesc(file_tesseral_backend_v1_audit_log_events, 28);

/**
 * @generated from message tesseral.backend.v1.UserRoleAssignmentCreated
 */
export type UserRoleAssignmentCreated = Message<"tesseral.backend.v1.UserRoleAssignmentCreated"> & {
  /**
   * @generated from field: tesseral.backend.v1.UserRoleAssignment user_role_assignment = 1;
   */
  userRoleAssignment?: UserRoleAssignment;
};

/**
 * Describes the message tesseral.backend.v1.UserRoleAssignmentCreated.
 * Use `create(UserRoleAssignmentCreatedSchema)` to create a new message.
 */
export const UserRoleAssignmentCreatedSchema: GenMessage<UserRoleAssignmentCreated> = /*@__PURE__*/
  messageDesc(file_tesseral_backend_v1_audit_log_events, 29);

/**
 * @generated from message tesseral.backend.v1.UserRoleAssignmentDeleted
 */
export type UserRoleAssignmentDeleted = Message<"tesseral.backend.v1.UserRoleAssignmentDeleted"> & {
  /**
   * @generated from field: tesseral.backend.v1.UserRoleAssignment user_role_assignment = 1;
   */
  userRoleAssignment?: UserRoleAssignment;
};

/**
 * Describes the message tesseral.backend.v1.UserRoleAssignmentDeleted.
 * Use `create(UserRoleAssignmentDeletedSchema)` to create a new message.
 */
export const UserRoleAssignmentDeletedSchema: GenMessage<UserRoleAssignmentDeleted> = /*@__PURE__*/
  messageDesc(file_tesseral_backend_v1_audit_log_events, 30);

/**
 * @generated from message tesseral.backend.v1.UserCreated
 */
export type UserCreated = Message<"tesseral.backend.v1.UserCreated"> & {
  /**
   * @generated from field: tesseral.backend.v1.User user = 1;
   */
  user?: User;
};

/**
 * Describes the message tesseral.backend.v1.UserCreated.
 * Use `create(UserCreatedSchema)` to create a new message.
 */
export const UserCreatedSchema: GenMessage<UserCreated> = /*@__PURE__*/
  messageDesc(file_tesseral_backend_v1_audit_log_events, 31);

/**
 * @generated from message tesseral.backend.v1.UserDeleted
 */
export type UserDeleted = Message<"tesseral.backend.v1.UserDeleted"> & {
  /**
   * @generated from field: tesseral.backend.v1.User user = 1;
   */
  user?: User;
};

/**
 * Describes the message tesseral.backend.v1.UserDeleted.
 * Use `create(UserDeletedSchema)` to create a new message.
 */
export const UserDeletedSchema: GenMessage<UserDeleted> = /*@__PURE__*/
  messageDesc(file_tesseral_backend_v1_audit_log_events, 32);

/**
 * @generated from message tesseral.backend.v1.UserUpdated
 */
export type UserUpdated = Message<"tesseral.backend.v1.UserUpdated"> & {
  /**
   * @generated from field: tesseral.backend.v1.User user = 1;
   */
  user?: User;

  /**
   * @generated from field: tesseral.backend.v1.User previous_user = 2;
   */
  previousUser?: User;
};

/**
 * Describes the message tesseral.backend.v1.UserUpdated.
 * Use `create(UserUpdatedSchema)` to create a new message.
 */
export const UserUpdatedSchema: GenMessage<UserUpdated> = /*@__PURE__*/
  messageDesc(file_tesseral_backend_v1_audit_log_events, 33);

