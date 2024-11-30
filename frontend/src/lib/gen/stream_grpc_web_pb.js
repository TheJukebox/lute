/**
 * @fileoverview gRPC-Web generated client stub for stream
 * @enhanceable
 * @public
 */

// Code generated by protoc-gen-grpc-web. DO NOT EDIT.
// versions:
// 	protoc-gen-grpc-web v1.5.0
// 	protoc              v3.6.1
// source: stream.proto


/* eslint-disable */
// @ts-nocheck



import {GrpcWebClientBase, MethodDescriptor, MethodType} from 'grpc-web'; 
import './stream_pb';

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?grpc.web.ClientOptions} options
 * @constructor
 * @struct
 * @final
 */
proto.stream.AudioStreamClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options.format = 'text';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname.replace(/\/+$/, '');

};


/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?grpc.web.ClientOptions} options
 * @constructor
 * @struct
 * @final
 */
proto.stream.AudioStreamPromiseClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options.format = 'text';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname.replace(/\/+$/, '');

};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.stream.AudioStreamRequest,
 *   !proto.stream.AudioStreamChunk>}
 */
const methodDescriptor_AudioStream_StreamAudio = new MethodDescriptor(
  '/stream.AudioStream/StreamAudio',
  MethodType.SERVER_STREAMING,
  proto.stream.AudioStreamRequest,
  proto.stream.AudioStreamChunk,
  /**
   * @param {!proto.stream.AudioStreamRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.stream.AudioStreamChunk.deserializeBinary
);


/**
 * @param {!proto.stream.AudioStreamRequest} request The request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!grpc.web.ClientReadableStream<!proto.stream.AudioStreamChunk>}
 *     The XHR Node Readable Stream
 */
proto.stream.AudioStreamClient.prototype.streamAudio =
    function(request, metadata) {
  return this.client_.serverStreaming(this.hostname_ +
      '/stream.AudioStream/StreamAudio',
      request,
      metadata || {},
      methodDescriptor_AudioStream_StreamAudio);
};


/**
 * @param {!proto.stream.AudioStreamRequest} request The request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!grpc.web.ClientReadableStream<!proto.stream.AudioStreamChunk>}
 *     The XHR Node Readable Stream
 */
proto.stream.AudioStreamPromiseClient.prototype.streamAudio =
    function(request, metadata) {
  return this.client_.serverStreaming(this.hostname_ +
      '/stream.AudioStream/StreamAudio',
      request,
      metadata || {},
      methodDescriptor_AudioStream_StreamAudio);
};


export default proto.stream;

