/*
 *
 * Copyright 2023 puzzleweaver authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package service

import "context"

type SessionService interface {
	SettingsService
	Generate(ctx context.Context) (uint64, error)
}

type TemplateService interface {
	Render(ctx context.Context, templateName string, data any) ([]byte, error)
}

type SettingsService interface {
	Get(ctx context.Context, id uint64) (map[string]string, error)
	Update(ctx context.Context, id uint64, info map[string]string) error
}

type PasswordStrengthService interface {
	Validate(ctx context.Context, password string) (bool, error)
	GetRules(ctx context.Context, lang string) (string, error)
}

type MarkdownService interface {
	Apply(ctx context.Context, text string) (string, error)
}
