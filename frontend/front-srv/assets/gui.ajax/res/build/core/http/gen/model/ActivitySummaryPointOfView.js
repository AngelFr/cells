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

"use strict";

exports.__esModule = true;

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { "default": obj }; }

function _classCallCheck(instance, Constructor) { if (!(instance instanceof Constructor)) { throw new TypeError("Cannot call a class as a function"); } }

var _ApiClient = require('../ApiClient');

var _ApiClient2 = _interopRequireDefault(_ApiClient);

/**
* Enum class ActivitySummaryPointOfView.
* @enum {}
* @readonly
*/

var ActivitySummaryPointOfView = (function () {
    function ActivitySummaryPointOfView() {
        _classCallCheck(this, ActivitySummaryPointOfView);

        this.GENERIC = "GENERIC";
        this.ACTOR = "ACTOR";
        this.SUBJECT = "SUBJECT";
    }

    /**
    * Returns a <code>ActivitySummaryPointOfView</code> enum value from a Javascript object name.
    * @param {Object} data The plain JavaScript object containing the name of the enum value.
    * @return {module:model/ActivitySummaryPointOfView} The enum <code>ActivitySummaryPointOfView</code> value.
    */

    ActivitySummaryPointOfView.constructFromObject = function constructFromObject(object) {
        return object;
    };

    return ActivitySummaryPointOfView;
})();

exports["default"] = ActivitySummaryPointOfView;
module.exports = exports["default"];

/**
 * value: "GENERIC"
 * @const
 */

/**
 * value: "ACTOR"
 * @const
 */

/**
 * value: "SUBJECT"
 * @const
 */
