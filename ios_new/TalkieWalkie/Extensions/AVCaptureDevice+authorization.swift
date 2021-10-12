//
//  AVCaptureDevice+authorization.swift
//  TalkieWalkie
//
//  Created by Alexandre Carlier on 05.10.21.
//

import AVFoundation

extension AVCaptureDevice {
    enum AuthorizationStatus {
        case justDenied
        case alreadyDenied
        case restricted
        case justAuthorized
        case alreadyAuthorized
        case unknown

        var superStatus: AVAuthorizationStatus {
            switch self {
            case .justDenied, .alreadyDenied:
                return .denied
            case .restricted:
                return .restricted
            case .justAuthorized, .alreadyAuthorized:
                return .authorized
            default:
                return .notDetermined
            }
        }
    }

    class func authorizeVideo(completion: ((AuthorizationStatus) -> Void)?) {
        AVCaptureDevice.authorize(mediaType: AVMediaType.video, completion: completion)
    }

    class func authorizeAudio(completion: ((AuthorizationStatus) -> Void)?) {
        AVCaptureDevice.authorize(mediaType: AVMediaType.audio, completion: completion)
    }

    private class func authorize(mediaType: AVMediaType, completion: ((AuthorizationStatus) -> Void)?) {
        let status = AVCaptureDevice.authorizationStatus(for: mediaType)
        switch status {
        case .authorized:
            completion?(.alreadyAuthorized)
        case .denied:
            completion?(.alreadyDenied)
        case .restricted:
            completion?(.restricted)
        case .notDetermined:
            AVCaptureDevice.requestAccess(for: mediaType, completionHandler: { granted in
                DispatchQueue.main.async {
                    if granted {
                        completion?(.justAuthorized)
                    } else {
                        completion?(.justDenied)
                    }
                }
            })
        @unknown default:
            completion?(.unknown)
        }
    }
}
