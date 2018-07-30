/**
 * Pydio Cells Rest API
 * No description provided (generated by Swagger Codegen https://github.com/swagger-api/swagger-codegen)
 *
 * OpenAPI spec version: 1.0
 * 
 *
 * NOTE: This class is auto generated by the swagger code generator program.
 * https://github.com/swagger-api/swagger-codegen.git
 * Do not edit the class manually.
 *
 */

'use strict';

exports.__esModule = true;

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { 'default': obj }; }

function _classCallCheck(instance, Constructor) { if (!(instance instanceof Constructor)) { throw new TypeError('Cannot call a class as a function'); } }

var _ApiClient = require('../ApiClient');

var _ApiClient2 = _interopRequireDefault(_ApiClient);

/**
* The RestUserJobResponse model module.
* @module model/RestUserJobResponse
* @version 1.0
*/

var RestUserJobResponse = (function () {
    /**
    * Constructs a new <code>RestUserJobResponse</code>.
    * @alias module:model/RestUserJobResponse
    * @class
    */

    function RestUserJobResponse() {
        _classCallCheck(this, RestUserJobResponse);

        this.JobUuid = undefined;
    }

    /**
    * Constructs a <code>RestUserJobResponse</code> from a plain JavaScript object, optionally creating a new instance.
    * Copies all relevant properties from <code>data</code> to <code>obj</code> if supplied or a new instance if not.
    * @param {Object} data The plain JavaScript object bearing properties of interest.
    * @param {module:model/RestUserJobResponse} obj Optional instance to populate.
    * @return {module:model/RestUserJobResponse} The populated <code>RestUserJobResponse</code> instance.
    */

    RestUserJobResponse.constructFromObject = function constructFromObject(data, obj) {
        if (data) {
            obj = obj || new RestUserJobResponse();

            if (data.hasOwnProperty('JobUuid')) {
                obj['JobUuid'] = _ApiClient2['default'].convertToType(data['JobUuid'], 'String');
            }
        }
        return obj;
    };

    /**
    * @member {String} JobUuid
    */
    return RestUserJobResponse;
})();

exports['default'] = RestUserJobResponse;
module.exports = exports['default'];