package services

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// OGPInfo はOGP情報を格納する構造体
type OGPInfo struct {
	Title       string
	Description string
	Image       string
	URL         string
}

// OGPService はOGP情報を取得するサービスのインターフェース
type OGPService interface {
	FetchOGP(ctx context.Context, targetURL string) (*OGPInfo, error)
}

type ogpService struct {
	httpClient *http.Client
}

// NewOGPService は新しいOGPServiceを作成します
func NewOGPService() OGPService {
	return &ogpService{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// FetchOGP は指定されたURLからOGP情報を取得します
func (s *ogpService) FetchOGP(ctx context.Context, targetURL string) (*OGPInfo, error) {
	// URLのバリデーション
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		return nil, fmt.Errorf("無効なURLです: %w", err)
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return nil, errors.New("URLはhttpまたはhttpsで始まる必要があります")
	}

	// HTTPリクエストを作成
	req, err := http.NewRequestWithContext(ctx, "GET", targetURL, nil)
	if err != nil {
		return nil, fmt.Errorf("リクエストの作成に失敗しました: %w", err)
	}

	// User-Agentを設定（一部のサイトでブロックされないように）
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; OGPBot/1.0)")

	// HTTPリクエストを実行
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ページの取得に失敗しました: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ページの取得に失敗しました: ステータスコード %d", resp.StatusCode)
	}

	// HTMLをパース
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("HTMLのパースに失敗しました: %w", err)
	}

	ogp := &OGPInfo{
		URL: targetURL,
	}

	// OGP情報を抽出
	doc.Find("meta").Each(func(i int, sel *goquery.Selection) {
		property, _ := sel.Attr("property")
		content, _ := sel.Attr("content")

		switch property {
		case "og:title":
			ogp.Title = content
		case "og:description":
			ogp.Description = content
		case "og:image":
			ogp.Image = content
		case "og:url":
			if content != "" {
				ogp.URL = content
			}
		}
	})

	// OGPがない場合は通常のmetaタグから取得を試みる
	if ogp.Title == "" {
		ogp.Title = doc.Find("title").First().Text()
	}

	if ogp.Description == "" {
		doc.Find("meta[name='description']").Each(func(i int, sel *goquery.Selection) {
			if content, exists := sel.Attr("content"); exists {
				ogp.Description = content
			}
		})
	}

	// タイトルが取得できなかった場合はエラー
	if ogp.Title == "" {
		return nil, errors.New("ページからタイトルを取得できませんでした")
	}

	// タイトルとディスクリプションをトリム
	ogp.Title = strings.TrimSpace(ogp.Title)
	ogp.Description = strings.TrimSpace(ogp.Description)

	return ogp, nil
}
