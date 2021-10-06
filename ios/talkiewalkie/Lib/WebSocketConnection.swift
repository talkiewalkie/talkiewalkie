//
//  WebSocketConnection.swift
//  talkiewalkie
//
//  Created by Théo Matussière on 01/10/2021.
//

import Combine
import Foundation
import OSLog

// Adapted from: https://github.dev/appspector/URLSessionWebSocketTask/tree/master/WebSockets

protocol WebSocketConnection {
    func send(text: String)
    func send(data: Data)
    func connect()
    func disconnect()
    var delegate: WebSocketConnectionDelegate? {
        get
        set
    }
}

protocol WebSocketConnectionDelegate {
    func onConnected(connection: WebSocketConnection)
    func onDisconnected(connection: WebSocketConnection, error: Error?)
    func onError(connection: WebSocketConnection, error: Error)
    func onMessage(connection: WebSocketConnection, text: String)
    func onMessage(connection: WebSocketConnection, data: Data)
}

class WebSocketTaskConnection: NSObject, WebSocketConnection, URLSessionWebSocketDelegate {
    var connected = false
    var numReconnects = 0

    var url: String { webSocketTask.currentRequest?.url?.absoluteString ?? "[no url yet]" }
    private let logger = Logger(subsystem: Bundle.main.bundleIdentifier!, category: "Network")
    private let pingIntervalSeconds = 10.0

    var delegate: WebSocketConnectionDelegate?
    var webSocketTask: URLSessionWebSocketTask
    var urlSession: URLSession
    var initialRequest: URLRequest
    let delegateQueue = OperationQueue()

    init(request: URLRequest) {
        urlSession = URLSession.shared
        initialRequest = request
        webSocketTask = urlSession.webSocketTask(with: request)
        super.init()
    }

    func urlSession(_: URLSession, webSocketTask _: URLSessionWebSocketTask, didOpenWithProtocol _: String?) {
        connected = true
        delegate?.onConnected(connection: self)
    }

    func urlSession(_: URLSession, webSocketTask _: URLSessionWebSocketTask, didCloseWith _: URLSessionWebSocketTask.CloseCode, reason _: Data?) {
        connected = false
        delegate?.onDisconnected(connection: self, error: nil)
    }

    func connect() {
        logger.debug("[\(self.url)] connecting to ws...")
        webSocketTask = urlSession.webSocketTask(with: initialRequest)
        webSocketTask.resume()

        // if we fail n times we'll have n loops sending pings at various moments...
        if numReconnects == 0 {
            DispatchQueue.main.async { self.keepAlive() }
        }
        listen()
    }

    func keepAlive() {
        webSocketTask.sendPing { error in
            if let error = error {
                self.logger.error("[\(self.url)] ping error: \(error.localizedDescription)")
                self.connected = false
            } else {
                self.logger.debug("[\(self.url)] ping success")
                self.connected = true
            }
            if self.connected {
                self.logger.debug("[\(self.url)] scheduling next ping")
                DispatchQueue.main.asyncAfter(deadline: .now() + self.pingIntervalSeconds) {
                    self.keepAlive()
                }
            }
        }
    }

    func disconnect() {
        logger.debug("[\(self.url)] disconnecting ws...")
        webSocketTask.cancel(with: .goingAway, reason: nil)
    }

    func listen() {
        webSocketTask.receive { result in
            switch result {
            case let .failure(error):
                self.delegate?.onError(connection: self, error: error)
            case let .success(message):
                switch message {
                case let .string(text):
                    self.delegate?.onMessage(connection: self, text: text)
                case let .data(data):
                    self.delegate?.onMessage(connection: self, data: data)
                @unknown default:
                    fatalError()
                }

                self.listen()
            }
        }
    }

    func send(text: String) {
        webSocketTask.send(URLSessionWebSocketTask.Message.string(text)) { error in
            if let error = error {
                self.delegate?.onError(connection: self, error: error)
            }
        }
    }

    func send(data: Data) {
        webSocketTask.send(URLSessionWebSocketTask.Message.data(data)) { error in
            if let error = error {
                self.delegate?.onError(connection: self, error: error)
            }
        }
    }
}
