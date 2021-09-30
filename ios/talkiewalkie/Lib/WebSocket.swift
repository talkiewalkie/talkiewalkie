//
//  WebSocket.swift
//  talkiewalkie
//
//  Created by Théo Matussière on 28/09/2021.
//

import Foundation


class WebSocketManager: NSObject, ObservableObject {
    private static let KEEP_ALIVE_INTERVAL_SECONDS: Double = 5
    
    private var session: URLSession?
    private var task: URLSessionWebSocketTask?
    private var shouldTryToReconnect = true
    private var onReceive: ((Result<URLSessionWebSocketTask.Message, Error>) -> Void)?
    
    @Published var connected: Bool = false
    
    var url: String {
        return task?.currentRequest?.url?.absoluteString ?? "[null task - no url]"
    }
    
    init(with urlRequest: URLRequest, onReceive: @escaping (Result<URLSessionWebSocketTask.Message, Error>) -> Void) {
        super.init()
        self.session = URLSession(configuration: .default, delegate: self, delegateQueue: OperationQueue())
        self.task = session?.webSocketTask(with: urlRequest)
        self.onReceive = onReceive
    }
    
    func reinit() {
        print("\(Date()) reconnecting ws connection to [\(url)]...")
        disconnect()
        session = URLSession(configuration: .default, delegate: self, delegateQueue: OperationQueue())
        task = session?.webSocketTask(with: task!.currentRequest!)
        connect()
    }
    
    func keepAlive() {
        task?.sendPing { error in
            if let error = error {
                let nsError = error as NSError
                print("\(Date()) Error when sending PING to '\(self.url)':\(error)")
                if nsError.domain == NSURLErrorDomain, nsError.code == 57 {
                    self.connected = false
                }
            } else {
                self.connected = true
                print("\(Date()) WebSocket connection '\(self.url)' alive")
                DispatchQueue.global().asyncAfter(deadline: .now() + WebSocketManager.KEEP_ALIVE_INTERVAL_SECONDS) {
                    if self.connected {
                        self.keepAlive()
                    }
                }
            }
        }
    }
    
    public func connect() {
        print("\(Date()) WEBSOCKET '\(url)'")
        shouldTryToReconnect = true
        task?.resume()
    }

    public func disconnect() {
        task?.cancel(with: .goingAway, reason: nil)
        connected = false
        shouldTryToReconnect = false
    }
    
    public func receive() {
        task?.receive { result in
            self.onReceive?(result)
            if self.connected {
                self.receive()
            }
        }
    }
}

extension WebSocketManager: URLSessionWebSocketDelegate {
    func urlSession(_ session: URLSession, webSocketTask: URLSessionWebSocketTask, didOpenWithProtocol protocol: String?) {
        connected = true
        print("\(Date()) Connected to '\(url)'!")
        
        receive()
        keepAlive()
    }

    func urlSession(_ session: URLSession, webSocketTask: URLSessionWebSocketTask, didCloseWith closeCode: URLSessionWebSocketTask.CloseCode, reason: Data?) {
        connected = false
        print("\(Date()) Disconnected from '\(url)' with close code: \(closeCode) - reason \(reason)!")
        if shouldTryToReconnect {
            print("\(Date()) Reconnecting...")
            connect()
        }
    }
    
    func urlSession(_ session: URLSession, task: URLSessionTask, didCompleteWithError: Error?) {
        if let err = didCompleteWithError {
            print("\(Date()) failed to connect to '\(url)': \(err)")
        }
    }
}
